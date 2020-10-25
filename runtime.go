package mesh

import (
	"sync"
)

// Type alias.
type ipAddress = int32

// Runtime for all p2p connections.
type P2pRuntime struct {
	// Way to spawn new connections.
	spawner ConnectionSpawner
	// Links ip to it's connections.
	ip2conn map[ipAddress]Connection
	// Keep sync.
	mu sync.Mutex
}

// Creates new P2pRuntime.
func NewP2pRuntime(spawner ConnectionSpawner) *P2pRuntime {
	return &P2pRuntime{
		ip2conn: make(map[ipAddress]Connection),
		spawner: spawner,
	}
}

// Returns `Connection` for ip address.
func (a *P2pRuntime) GetConnection(ip ipAddress) Connection {
	a.mu.Lock()
	defer a.mu.Unlock()
	if conn, ok := a.ip2conn[ip]; ok {
		return conn
	}
	connected := a.spawner.SpawnOutgoing(ip)
	a.ip2conn[ip] = connected
	return connected
}

// Called when new connection arrived.
func (a *P2pRuntime) OnNewRemoteConnection(ip ipAddress, conn Connection) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.ip2conn[ip]; ok {
		// Seem like we already have same connection, so just close duplicate.
		conn.Close()
		return
	}
	a.ip2conn[ip] = conn
}

// Closes all connections.
func (a *P2pRuntime) Shutdown() {
	a.mu.Lock()
	for _, conn := range a.ip2conn {
		conn.Close()
	}
	a.ip2conn = make(map[ipAddress]Connection)
	a.mu.Unlock()
}

// Returns number of alive connections.
func (a *P2pRuntime) Len() int {
	a.mu.Lock()
	cnt := len(a.ip2conn)
	a.mu.Unlock()
	return cnt
}
