package mesh

// Interface that help us to spawn outgoing connections.
// Main idea that when we spawn outgoing connection, we provide only ip address,
// but mostly we need protocol version, node name, and other parameters, so this will
// stored inside.
type ConnectionSpawner interface {
	SpawnOutgoing(ip int32) Connection
}
