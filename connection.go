package mesh

// Interface for network connection.
type Connection interface {
	Close()
	Closed() bool
	Open()
}

type ConnectionImpl struct {
	closed bool
}

func (a *ConnectionImpl) Close() {
	a.closed = true
}

func (a ConnectionImpl) Closed() bool {
	return a.closed
}

func (ConnectionImpl) Open() {}

func NewConnection() Connection {
	return &ConnectionImpl{}
}

type ClosedConnection struct {
}

func (a *ClosedConnection) Close() {
}

func (a ClosedConnection) Closed() bool {
	return true
}

func (a ClosedConnection) Open() {
}
