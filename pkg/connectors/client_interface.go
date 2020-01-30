package connectors

import (
	gocb "github.com/couchbase/gocb/v2"
)

// Client Interface - used as a receiver and can be overriden for testing
type Clients interface {
	AnalyticsQuery(query string, opts *gocb.AnalyticsOptions) (*gocb.AnalyticsResult, error)
	Close() error
}
