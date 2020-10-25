package mesh

// Interface for network connection.
type Connection interface {
	Close()
	Closed() bool
}
