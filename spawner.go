package mesh

import (
	"sync"
	"sync/atomic"
)

var skippedReceivers = uint64(0)

// Interface that help us to spawn outgoing connections.
// Main idea that when we spawn outgoing connection, we provide only ip address,
// but mostly we need protocol version, node name, and other parameters, so this will
// stored inside.
type ConnectionSpawner interface {
	SpawnOutgoing(ip int32) chan Connection
	Close()
}

type ConnectionSpawnerImpl struct {
	mu                sync.Mutex
	times             map[ipAddress]int
	connectionFactory func() Connection
	closed            int32
}

func NewConnectionSpawner(connectionFactory func() Connection) *ConnectionSpawnerImpl {
	return &ConnectionSpawnerImpl{
		times:             make(map[ipAddress]int),
		connectionFactory: connectionFactory,
	}
}

func (a *ConnectionSpawnerImpl) SpawnOutgoing(ip int32) chan Connection {
	a.mu.Lock()
	defer a.mu.Unlock()

	ch := make(chan Connection, 1)
	// No real actions when we closed.
	if atomic.LoadInt32(&a.closed) == 1 {
		ch <- &ClosedConnection{}
		return ch
	}

	if a.times[ip] == 0 {
		go func() {
			c := a.connectionFactory()
			c.Open()
			if atomic.LoadInt32(&a.closed) == 1 {
				c.Close()
				return
			}
			a.mu.Lock()
			times := a.times[ip]
			for i := 0; i <= times; i++ {
				select {
				case ch <- c:
				//We need default cause if waiter failed somehow, we should not hung forever.
				default:
					atomic.AddUint64(&skippedReceivers, 1)
				}
			}
			delete(a.times, ip)
			a.mu.Unlock()
		}()
	} else {
		a.times[ip]++
	}

	return ch
}

func (a *ConnectionSpawnerImpl) Close() {
	atomic.StoreInt32(&a.closed, 1)
}
