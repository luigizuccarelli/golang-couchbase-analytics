package connectors

import (
	"fmt"
	"os"

	gocb "github.com/couchbase/gocb/v2"
	"github.com/microlib/simple"
)

// Connectors struct - all backend connections in a common object
type Connectors struct {
	Bucket  *gocb.Bucket
	Cluster *gocb.Cluster
	Logger  *simple.Logger
}

type Result struct {
	Data *gocb.AnalyticsResult
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

// Upsert : wrapper function for couchbase update
func (r Result) Row(ptr interface{}) error {
	return r.Data.Row(ptr)
}

func (r Result) Next() bool {
	return r.Data.Next()
}

func (r Result) Close() error {
	return r.Data.Close()
}

func (conn Connectors) AnalyticsQuery(query string, opts *gocb.AnalyticsOptions) (AnalyticsResult, error) {
	res, err := conn.Cluster.AnalyticsQuery(query, opts)
	//ar := NewResult(res)
	return res, err
}

func (conn Connectors) Close() {
	conn.Cluster.Close(nil)
}

// NewClientConnectors returns Connectors struct
func NewClientConnections(logger *simple.Logger) Clients {

	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: os.Getenv("COUCHBASE_USER"),
			Password: os.Getenv("COUCHBASE_PASSWORD"),
		},
	}

	cluster, err := gocb.Connect(os.Getenv("COUCHBASE_HOST"), opts)
	if err != nil {
		panic(err)
	}

	// get a bucket reference
	bucket := cluster.Bucket(os.Getenv("COUCHBASE_BUCKET"))

	conns := &Connectors{Bucket: bucket, Cluster: cluster, Logger: logger}
	return conns
}
