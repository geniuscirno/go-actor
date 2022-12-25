package core

import (
	"context"
)

type Node interface {
	Name() string
	Spawn(behavior ProcessBehavior, opts *SpawnOptions) (Process, error)
	SendMessage(ctx context.Context, to PID, message Message) error
	Stop()
	Wait()
	Join(c Cluster)
}

type node struct {
	ctx    context.Context
	cancel context.CancelFunc

	name string

	registry *ProcessRegistry

	cluster Cluster
}

func NewNode(name string) Node {
	ctx, cancel := context.WithCancel(context.Background())
	node := &node{
		ctx:    ctx,
		cancel: cancel,

		name:     name,
		registry: newProcessRegistry(),
	}
	return node
}

func (n *node) newPID() PID {
	return PID{
		Node: n.name,
		ID:   n.registry.NextId(),
	}
}

func (n *node) newProcess(parent *process, behavior ProcessBehavior, opts *SpawnOptions) (*process, error) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	if parent != nil {
		ctx, cancel = context.WithCancel(parent.ctx)
	} else {
		ctx, cancel = context.WithCancel(n.ctx)
	}

	p := &process{
		ctx:    ctx,
		cancel: cancel,

		node: n,

		behavior: behavior,

		mailbox: make(chan Message, 10000),
		exit:    make(chan struct{}, 1),

		parent: parent,
	}

	if opts.Name == "" {
		p.pid = n.newPID()
	} else {
		p.pid = PID{Node: n.name, ID: opts.Name}
	}

	if err := n.registry.Add(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (n *node) spawn(parent *process, behavior ProcessBehavior, opts *SpawnOptions) (*process, error) {
	p, err := n.newProcess(parent, behavior, opts)
	if err != nil {
		return nil, err
	}

	if parent != nil {
		parent.addChild(p)
	}
	cleanProcess := func(err error) {
		if parent != nil {
			parent.deleteChild(p)
		}
		n.registry.Delete(p)

		p.cancel()
	}

	go func(p Process) {
		err := behavior.ProcessLoop(p)
		cleanProcess(err)
	}(p)
	return p, nil
}

func (n *node) Name() string {
	return n.name
}

func (n *node) Spawn(behavior ProcessBehavior, opts *SpawnOptions) (Process, error) {
	return n.spawn(nil, behavior, opts)
}

func (n *node) SendMessage(ctx context.Context, to PID, message Message) error {
	if to.Node != n.name {
		return n.cluster.SendMessage(ctx, to, message)
	}

	process, err := n.registry.Get(to)
	if err != nil {
		return err
	}

	select {
	case process.mailbox <- message:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (n *node) Stop() {
	n.cancel()
}

func (n *node) Wait() {
	<-n.ctx.Done()
}

func (n *node) Join(c Cluster) {
	n.cluster = c
}
