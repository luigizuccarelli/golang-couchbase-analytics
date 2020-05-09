package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	//"reflect"
	"testing"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/trackmate-couchbase-analytics/pkg/connectors"
	gocb "github.com/couchbase/gocb/v2"
	"github.com/microlib/simple"
)

type FakeCouchbase struct {
}

type FakeResult struct {
	Type          string
	QueryError    bool
	FunctionError bool
}

func NewFakeResult(object string, queryErr bool, functionErr bool) connectors.AnalyticsResult {
	return FakeResult{
		Type:          object,
		QueryError:    queryErr,
		FunctionError: functionErr,
	}
}

func (conn Connectors) AnalyticsQuery(query string, opts *gocb.AnalyticsOptions) (connectors.AnalyticsResult, error) {
	if conn.QueryError {
		return NewFakeResult(conn.Type, true, false), errors.New("Fake error")
	}
	return NewFakeResult(conn.Type, false, conn.FunctionError), nil
}

func (conn Connectors) Close() {
}

func (c *Connectors) Error(msg string, val ...interface{}) {
	c.Logger.Error(fmt.Sprintf(msg, val...))
}

func (c *Connectors) Info(msg string, val ...interface{}) {
	c.Logger.Info(fmt.Sprintf(msg, val...))
}

func (c *Connectors) Debug(msg string, val ...interface{}) {
	c.Logger.Debug(fmt.Sprintf(msg, val...))
}

func (c *Connectors) Trace(msg string, val ...interface{}) {
	c.Logger.Trace(fmt.Sprintf(msg, val...))
}

func (fs FakeResult) Next(data interface{}) bool {
	m := make(map[string]interface{})
	if fs.Type == "Sankey" {
		m["count"] = "10"
		m["source"] = "AB"
		m["destination"] = "XY"
		*data.(*map[string]interface{}) = m
	}

	if fs.Type == "Funnel" {
		m["count"] = "10"
		m["pagename"] = "Test"
		m["pagetype"] = "Landing"
		m["utm_source"] = "mail"
		*data.(*map[string]interface{}) = m
	}
	return false
}

func (fs FakeResult) Close() error {
	if fs.FunctionError {
		return errors.New("Fake connection close")
	}
	return nil
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewHttpTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

type Connectors struct {
	Http          *http.Client
	Bucket        *FakeCouchbase
	Logger        *simple.Logger
	Type          string
	QueryError    bool
	FunctionError bool
}

type Parameters struct {
	File          string
	Code          int
	Type          string
	QueryError    bool
	FunctionError bool
}

func NewTestConnectors(params *Parameters, logger *simple.Logger) connectors.Clients {

	// we first load the json payload to simulate a call to middleware
	// for now just ignore failures.
	data, err := ioutil.ReadFile(params.File)
	if err != nil {
		logger.Error(fmt.Sprintf("file data %v\n", err))
		panic(err)
	}
	httpclient := NewHttpTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: params.Code,
			// Send response to be tested

			Body: ioutil.NopCloser(bytes.NewBufferString(string(data))),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	// conns := &Connectors{Http: httpclient, Bucket: &FakeCouchbase{}, Cluster: NewFakeCluster()}
	conns := &Connectors{Http: httpclient, Bucket: &FakeCouchbase{}, Type: params.Type, QueryError: params.QueryError, FunctionError: params.FunctionError, Logger: logger}
	return conns
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestAllMiddleware(t *testing.T) {

	logger := &simple.Logger{Level: "trace"}

	t.Run("IsAlive : should pass", func(t *testing.T) {
		var STATUS int = 200
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v2/sys/info/isalive", nil)
		p := &Parameters{File: "../../tests/payload.json", Code: STATUS, Type: "NA", QueryError: false, FunctionError: false}
		conn := NewTestConnectors(p, logger)
		handler := http.HandlerFunc(IsAlive)
		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		conn.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "IsAlive", rr.Code, STATUS))
		}
	})

	t.Run("SankeyChartHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		data, _ := ioutil.ReadFile("../../tests/payload.json")
		req, _ := http.NewRequest("POST", "/api/v1/sankeydata", bytes.NewBuffer(data))
		p := &Parameters{File: "../../tests/payload.json", Code: STATUS, Type: "Sankey", QueryError: false, FunctionError: false}
		conn := NewTestConnectors(p, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SankeyChartHandler(w, r, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "SankeyChartHandler", rr.Code, STATUS))
		}
	})

	t.Run("SankeyChartHandler : should fail", func(t *testing.T) {
		var STATUS int = 500
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		data, _ := ioutil.ReadFile("../../tests/payload.json")
		req, _ := http.NewRequest("POST", "/api/v1/sankeydata", bytes.NewBuffer(data))
		p := &Parameters{File: "../../tests/payload.json", Code: STATUS, Type: "Sankey", QueryError: true, FunctionError: false}
		conn := NewTestConnectors(p, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SankeyChartHandler(w, r, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "SankeyChartHandler", rr.Code, STATUS))
		}
	})

	t.Run("SankeyChartHandler : should fail", func(t *testing.T) {
		var STATUS int = 500
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		data, _ := ioutil.ReadFile("../../tests/payload.json")
		req, _ := http.NewRequest("POST", "/api/v1/sankeydata", bytes.NewBuffer(data))
		p := &Parameters{File: "../../tests/payload.json", Code: STATUS, Type: "Sankey", QueryError: false, FunctionError: true}
		conn := NewTestConnectors(p, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SankeyChartHandler(w, r, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "SankeyChartHandler", rr.Code, STATUS))
		}
	})
}
