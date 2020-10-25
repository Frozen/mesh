package mesh

import "testing"

type connectionMock struct {
	closed bool
}

// Close connection.
func (a *connectionMock) Close() {
	a.closed = true
}

// Returns `true` if connection closed.
func (a *connectionMock) Closed() bool {
	return a.closed
}

// Just check it compiles.
func TestConnection_Smoke(t *testing.T) {
	var _ Connection = &connectionMock{}
}
