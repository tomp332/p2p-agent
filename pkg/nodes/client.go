package nodes

type NodeClient interface {
	Connect() (string, error)
	Disconnect() error
}
