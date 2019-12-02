package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/couchbase/gocb.v1"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

func SankeyChartHandler(w http.ResponseWriter, r *http.Request) {
	var response Response

	addHeaders(w, r)

	q := "select `from`.`pagename` as source,`to`.`pagename` as destination, count(`trackingid`) as count  from SBR where `event`.`type` = 'load'  group by `from`.`pagename` as source, `to`.`pagename` as destination"
	query := gocb.NewAnalyticsQuery(q)

	results, err := cluster.ExecuteAnalyticsQuery(query, nil)
	if err != nil {
		response = Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf("Could not execute analytics query from couchbase %v", err)}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// all good :)
	d, err := processSankeyResults(results)
	if err != nil {
		response = Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf("Could not process analytics results from couchbase %v", err)}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Debug(fmt.Sprintf("Analytics result from couchbase  %v \n", d))
	response = Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: "Data retrieved succesfully", Sankey: d}
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	logger.Debug(fmt.Sprintf("AnaytlicsHandler response : %s", string(b)))
	fmt.Fprintf(w, string(b))
}

func FunnelChartHandler(w http.ResponseWriter, r *http.Request) {
	var response Response
	var params []FunnelInputData

	addHeaders(w, r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf("Could not read body data %v\n", err)}
		w.WriteHeader(http.StatusInternalServerError)
	}

	// we first unmarshal the payload and add needed values before posting to couchbase
	errs := json.Unmarshal(body, &params)
	if errs != nil {
		logger.Error(fmt.Sprintf("Could not unmarshal analytics data to json %v", errs))
		response = Response{StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf("Could not read body data %v\n", err)}
	}

	var source, node string
	for x, _ := range params {
		if params[x].Type == "source" {
			source = source + "'" + params[x].Node + "',"
		}
		if params[x].Type == "node" {
			node = node + "'" + params[x].Node + "',"
		}
	}

	logger.Trace(fmt.Sprintf("LMZ DEBUG %v %s %s\n", params, source[:len(source)-1], node[:len(node)-1]))

	q := "select `utm_source`,`to`.`pagename` , `to`.`pagetype`, count(*) as count from SBR where `event`.`type` = 'load' and utm_campaign = 'WinBig' and utm_source in [ " + source[:len(source)-1] + " ] and `to`.`pagename` in [ " + node[:len(node)-1] + " ] group by `utm_source`,`to`.`pagename` , `to`.`pagetype`"
	logger.Trace(fmt.Sprintf("LMZ DEBUG %s\n", q))
	query := gocb.NewAnalyticsQuery(q)

	results, err := cluster.ExecuteAnalyticsQuery(query, nil)
	if err != nil {
		response = Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf("Could not execute analytics query from couchbase %v", err)}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// all good :)
	d, err := processFunnelResults(results)
	if err != nil {
		response = Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf("Could not process analytics results from couchbase %v", err)}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Debug(fmt.Sprintf("Analytics result from couchbase  %v \n", d))
	response = Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: "Data retrieved succesfully", Funnel: d}
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(response, "", "	")
	logger.Debug(fmt.Sprintf("AnaytlicsHandler response : %s", string(b)))
	fmt.Fprintf(w, string(b))
}

func SourceDropdownHandler(w http.ResponseWriter, r *http.Request) {
	var response Response
	vars := mux.Vars(r)

	logger.Debug(fmt.Sprintf("Mux vars %v\n", vars))

	addHeaders(w, r)

	q := "select distinct `utm_affiliate`,`utm_campaign`,`utm_source`,`utm_content` from SBR where `utm_affiliate` = '" + vars["affiliate"] + "' order by utm_source asc, utm_content asc"
	query := gocb.NewAnalyticsQuery(q)

	results, err := cluster.ExecuteAnalyticsQuery(query, nil)
	if err != nil {
		response = Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf("Could not execute query from couchbase (sourcedropdown) %v", err)}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// all good :)
	d, err := processSourceDropdownResults(results)
	if err != nil {
		response = Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf("Could not process results from couchbase (sourcedropdown) %v", err)}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Debug(fmt.Sprintf("SourceDropdown result from couchbase  %v \n", d))
	response = Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: "Data retrieved succesfully", SourceDropdown: d}
	w.WriteHeader(http.StatusOK)
	//}
	b, _ := json.MarshalIndent(response, "", "	")
	logger.Debug(fmt.Sprintf("SourceDropdownHandler response : %s", string(b)))
	fmt.Fprintf(w, string(b))
}

func DestinationDropdownHandler(w http.ResponseWriter, r *http.Request) {
	var response Response
	vars := mux.Vars(r)

	logger.Debug(fmt.Sprintf("Mux vars %v\n", vars))

	addHeaders(w, r)
	// TODO the from SBR must change - for each affiliate the name will be specific
	q := "select distinct `from`.`pagename` as source ,`to`.`pagename` as destination from SBR where `utm_campaign` = '" + vars["campaign"] + "' group by `from`.`pagename` as source, `to`.`pagename` as destination"
	query := gocb.NewAnalyticsQuery(q)

	results, err := cluster.ExecuteAnalyticsQuery(query, nil)
	if err != nil {
		response = Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf("Could not execute query from couchbase (destinationdropdown) %v", err)}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// all good :)
	d, err := processNodelinkResults(results)
	if err != nil {
		response = Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf("Could not process results from couchbase (destinationdropdown) %v", err)}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Debug(fmt.Sprintf("DestinationDropdown result from couchbase  %v \n", d))
	response = Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: "Data retrieved succesfully", DestinationDropdown: d}
	w.WriteHeader(http.StatusOK)
	//}
	b, _ := json.MarshalIndent(response, "", "	")
	logger.Debug(fmt.Sprintf("DestinationDropdownHandler response : %s", string(b)))
	fmt.Fprintf(w, string(b))
}

func NodelinkHandler(w http.ResponseWriter, r *http.Request) {
	var response Response
	vars := mux.Vars(r)

	logger.Debug(fmt.Sprintf("Mux vars %v\n", vars))

	addHeaders(w, r)
	// TODO the from SBR must change - for each affiliate the name will be specific
	q := "select distinct  `from`.`pagename` as source ,`to`.`pagename` as destination from SBR where `utm_affiliate` = '" + vars["affiliate"] + "' and `utm_campaign` = '" + vars["campaign"] + "' group by `from`.`pagename` as source, `to`.`pagename` as destination"
	query := gocb.NewAnalyticsQuery(q)

	results, err := cluster.ExecuteAnalyticsQuery(query, nil)
	if err != nil {
		response = Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf("Could not execute query from couchbase (nodelink) %v", err)}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// all good :)
	d, err := processNodelinkResults(results)
	if err != nil {
		response = Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf("Could not process results from couchbase (nodelink) %v", err)}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Debug(fmt.Sprintf("Nodelink result from couchbase  %v \n", d))
	response = Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: "Data retrieved succesfully", Nodelink: d}
	w.WriteHeader(http.StatusOK)
	//}
	b, _ := json.MarshalIndent(response, "", "	")
	logger.Debug(fmt.Sprintf("NodelinkHandler response : %s", string(b)))
	fmt.Fprintf(w, string(b))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ \"version\" : \"1.0.2\" , \"name\": \""+os.Getenv("NAME")+"\" }")
}

// headers (with cors) utility
func addHeaders(w http.ResponseWriter, r *http.Request) {
	var request []string
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	logger.Trace(fmt.Sprintf("Headers : %s", request))

	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

}

func processSankeyResults(results gocb.AnalyticsResults) ([]SankeyData, error) {
	var row map[string]interface{}
	var data []SankeyData

	for results.Next(&row) {
		if count, ok := row["count"]; ok {
			x, _ := strconv.Atoi(fmt.Sprintf("%v", count))
			d := SankeyData{Value: x, From: fmt.Sprintf("%v", row["source"]), To: fmt.Sprintf("%v", row["destination"])}
			data = append(data, d)
		}
	}
	if err := results.Close(); err != nil {
		return nil, err
	}

	return data, nil
}

func processFunnelResults(results gocb.AnalyticsResults) ([]FunnelOutputData, error) {
	var row map[string]interface{}
	var data []FunnelOutputData

	for results.Next(&row) {
		if count, ok := row["count"]; ok {
			x, _ := strconv.ParseInt(fmt.Sprintf("%v", count), 10, 64)
			d := FunnelOutputData{Count: x, PageName: fmt.Sprintf("%v", row["pagename"]), PageType: fmt.Sprintf("%v", row["pagetype"]), Source: fmt.Sprintf("%v", row["utm_source"])}
			data = append(data, d)
		}
	}
	if err := results.Close(); err != nil {
		return nil, err
	}

	return data, nil
}

func processSourceDropdownResults(results gocb.AnalyticsResults) ([]SourceDropdownData, error) {
	var row map[string]interface{}
	var data []SourceDropdownData

	for results.Next(&row) {
		if _, ok := row["utm_affiliate"]; ok {
			d := SourceDropdownData{Affiliate: fmt.Sprintf("%v", row["utm_affiliate"]),
				Campaign: fmt.Sprintf("%v", row["utm_campaign"]),
				Source:   fmt.Sprintf("%v", row["utm_source"]),
				Content:  fmt.Sprintf("%v", row["utm_content"])}

			data = append(data, d)
		}
	}
	if err := results.Close(); err != nil {
		return nil, err
	}

	return data, nil
}

func processNodelinkResults(results gocb.AnalyticsResults) ([]NodelinkData, error) {
	var row map[string]interface{}
	var data []NodelinkData

	for results.Next(&row) {
		if _, ok := row["source"]; ok {
			d := NodelinkData{Source: fmt.Sprintf("%v", row["source"]),
				Destination: fmt.Sprintf("%v", row["destination"])}

			data = append(data, d)
		}
	}
	if err := results.Close(); err != nil {
		return nil, err
	}

	return data, nil
}
