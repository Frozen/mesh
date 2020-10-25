package mesh

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRuntime(t *testing.T) {
	var (
		connIn  = &connectionMock{}
		spawner = &connectionSpawnerMock{
			timeout: 1 * time.Second,
		}
	)
	r := NewP2pRuntime(spawner)

	// Ensure no connection yet.
	require.Equal(t, 0, r.Len())

	outgoingConn := r.GetConnection(123)

	// Ensure we have only 1 connection.
	require.Equal(t, 1, r.Len())

	duplicateConn := r.GetConnection(123)

	// Again only 1 connection.
	require.Equal(t, 1, r.Len())

	// Because we return reference to interface, interface have same link to memory, so they equal.
	require.Equal(t, outgoingConn, duplicateConn)

	// Add incoming connection.
	r.OnNewRemoteConnection(321, connIn)

	// Ensure we have 2 different connections.
	require.Equal(t, 2, r.Len())

	// Add Duplicate incoming connection.
	r.OnNewRemoteConnection(123, &connectionMock{})

	// No new connection added, same quantity.
	require.Equal(t, 2, r.Len())

	// Close all connections.
	r.Shutdown()

	// Ensure no connection right now.
	require.Equal(t, 0, r.Len())

	// Ensure all connections closed
	require.True(t, connIn.Closed())
	require.True(t, outgoingConn.Closed())
}
