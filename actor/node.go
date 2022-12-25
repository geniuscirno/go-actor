package actor

import (
	"github.com/geniuscirno/go-actor/actor/future"
	"github.com/geniuscirno/go-actor/core"
)

type Node struct {
	core.Node
}

func NewNode(name string) *Node {
	return &Node{
		Node: core.NewNode(name),
	}
}

func (n *Node) SpawnActor(actor Actor, opt ...SpawnOption) (Process, error) {
	opts := &core.SpawnOptions{}
	for _, o := range opt {
		o(opts)
	}

	p, err := n.Spawn(&actorBehavior{
		futures: make(map[int64]*future.Future),
		actor:   actor,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &actorProcess{Process: p}, nil
}
