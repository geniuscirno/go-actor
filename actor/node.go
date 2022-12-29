package actor

import (
	"github.com/geniuscirno/go-actor/core"
)

type Node struct {
	core.Node
	Root Process
}

func NewNode(name string) *Node {
	n := &Node{
		Node: core.NewNode(name),
	}
	root, err := n.SpawnActor(ActorFunc(func(c Context) {}), Name("root"))
	if err != nil {
		panic(err)
	}
	n.Root = root
	return n
}

func (n *Node) SpawnActor(actor Actor, opt ...SpawnOption) (Process, error) {
	opts := &core.SpawnOptions{}
	for _, o := range opt {
		o(opts)
	}

	p, err := n.Spawn(&actorBehavior{
		actor: actor,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &actorProcess{Process: p}, nil
}
