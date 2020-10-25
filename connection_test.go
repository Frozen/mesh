package mesh

import (
	"testing"
	"time"
)

type connectionMock struct {
	closed  bool
	timeout time.Duration
}

// Close connection.
func (a *connectionMock) Close() {
	a.closed = true
}

// Returns `true` if connection closed.
func (a *connectionMock) Closed() bool {
	return a.closed
}

func (a connectionMock) Open() {
	<-time.After(a.timeout)
}

// Just check it compiles.
func TestConnection_Smoke(t *testing.T) {
	var _ Connection = &connectionMock{}
}
