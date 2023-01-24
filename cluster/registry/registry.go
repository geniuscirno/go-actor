package registry

import "context"

type Node struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Registrar interface {
	Register(ctx context.Context, node *Node) error
	Deregister(ctx context.Context, node *Node) error
	KeepAlive(ctx context.Context) error
}
