package mesh

import "time"

// Mock real spawner for testing purpose.
type connectionSpawnerMock struct {
	timeout    time.Duration
	connection Connection
}

// Mock real spawner for testing purpose.
func (a connectionSpawnerMock) SpawnOutgoing(ip int32) Connection {
	<-time.After(a.timeout)
	if a.connection == nil {
		return &connectionMock{}
	}
	return a.connection
}
