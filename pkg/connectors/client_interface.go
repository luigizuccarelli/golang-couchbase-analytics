package connectors

import (
	gocb "github.com/couchbase/gocb/v2"
)

type AnalyticsResult interface {
	Next(data interface{}) bool
	Close() error
}

type Clients interface {
	Error(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
	AnalyticsQuery(query string, opts *gocb.AnalyticsOptions) (AnalyticsResult, error)
	Close()
}
