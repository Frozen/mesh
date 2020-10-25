package mesh

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRuntime(t *testing.T) {
	var (
		connIn  = &connectionMock{}
		spawner = NewConnectionSpawner(func() Connection {
			return &connectionMock{
				timeout: 5 * time.Second,
			}
		})
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

// 3 connections should hang for 1 second, not 3 seconds.
func TestMultipleConnectionsSimultaneously(t *testing.T) {
	now := time.Now()
	var (
		spawner = NewConnectionSpawner(func() Connection {
			return &connectionMock{
				timeout: 1 * time.Second,
			}
		})
	)
	r := NewP2pRuntime(spawner)

	c1 := r.GetConnection(123)
	c2 := r.GetConnection(123)
	_ = r.GetConnection(123)

	require.True(t, c1 == c2)

	require.True(t, time.Now().Sub(now) < 2*time.Second)
}

// This code should work more than 2 seconds.
func TestSequentialConnections(t *testing.T) {
	now := time.Now()
	var (
		spawner = NewConnectionSpawner(func() Connection {
			return &connectionMock{
				timeout: 1 * time.Second,
			}
		})
	)
	r := NewP2pRuntime(spawner)

	c1 := r.GetConnection(123)
	c2 := r.GetConnection(321)

	require.False(t, c1 == c2)

	require.True(t, time.Now().Sub(now) > 2*time.Second)
}

func TestConnectionsOnShutdownRuntime(t *testing.T) {
	var (
		spawner = NewConnectionSpawner(func() Connection {
			return &connectionMock{
				timeout: 1 * time.Second,
			}
		})
	)
	r := NewP2pRuntime(spawner)
	r.Shutdown()
	c1 := r.GetConnection(123)

	require.True(t, c1.Closed())
}
