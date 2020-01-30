package connectors

import (
	"testing"

	gocb "github.com/couchbase/gocb/v2"
	"github.com/microlib/simple"
)

type FakeCouchbase struct {
}

type FakeCluster struct {
}

// Mock all connections
type MockConnections struct {
	Bucket  *FakeCouchbase
	Cluster *FakeCluster
}

func (r *MockConnections) Close() error {
	return nil
}

func (r *MockConnections) AnalyticsQuery(query string, opts *gocb.AnalyticsOptions) (*gocb.AnalyticsResult, error) {
	return &gocb.AnalyticsResult{}, nil
}

// NewTestConnections - create all mock connections
func NewTestConnections(file string, code int, logger *simple.Logger) Clients {

	bc := &FakeCouchbase{}
	conns := &MockConnections{Bucket: bc}
	return conns
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
