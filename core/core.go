package core

import (
	"errors"
)

var (
	ErrDupProcessName  = errors.New("dup process name")
	ErrProcessNotFound = errors.New("process not found")
	ErrProcessBusy     = errors.New("process busy")
	ErrTimeout         = errors.New("timeout")
)

//type core struct {
//	ctx    context.Context
//	cancel context.CancelFunc
//
//	node     *node
//	registry *ProcessRegistry
//}
//
//func newCore(node *node) *core {
//	ctx, cancel := context.WithCancel(context.Background())
//
//	c := &core{
//		ctx:    ctx,
//		cancel: cancel,
//		node:   node,
//	}
//	c.registry = newProcessRegistry()
//	return c
//}
//
//func (c *core) newPID() PID {
//	return PID{
//		node: c.node.name,
//		ID:   c.registry.NextId(),
//	}
//}
//
//func (c *core) newProcess(name string, behavior ProcessBehavior, parent *process) (*process, error) {
//	var (
//		ctx    context.Context
//		cancel context.CancelFunc
//	)
//
//	if parent != nil {
//		ctx, cancel = context.WithCancel(parent.ctx)
//	} else {
//		ctx, cancel = context.WithCancel(c.ctx)
//	}
//
//	p := &process{
//		ctx:    ctx,
//		cancel: cancel,
//
//		node: c.node,
//
//		pid:      c.newPID(),
//		behavior: behavior,
//
//		mailbox: make(chan Message, 100),
//
//		//futures: make(map[int64]*future.Future),
//
//		parent: parent,
//	}
//
//	if name == "" {
//		p.pid = c.newPID()
//	} else {
//		p.pid = PID{node: c.node.name, ID: name}
//	}
//
//	if err := c.registry.Add(p); err != nil {
//		return nil, err
//	}
//	return p, nil
//}

//func (c *core) spawn(name string, behavior ProcessBehavior, parent *process) (*process, error) {
//	p, err := c.newProcess(name, behavior, parent)
//	if err != nil {
//		return nil, err
//	}
//	if parent != nil {
//		parent.addChild(p)
//	}
//
//	cleanProcess := func(err error) {
//		if parent != nil {
//			parent.deleteChild(p)
//		}
//		c.registry.Delete(p)
//		p.cancel()
//	}
//
//	go func(p actorProcess) {
//		err := behavior.ProcessLoop(p)
//		cleanProcess(err)
//	}(p)
//	return p, nil
//}

//func (c *core) send(to PID, message Message) error {
//	process, err := c.registry.Get(to)
//	if err != nil {
//		return err
//	}
//
//	timer := time.NewTimer(time.Second * 30)
//	defer timer.Stop()
//
//	select {
//	case process.mailbox <- message:
//	case <-timer.C:
//		return ErrTimeout
//	}
//	return nil
//}
//
//func (c *core) Stop() {
//	c.cancel()
//}
//
//func (c *core) Wait() {
//	<-c.ctx.Done()
//}
