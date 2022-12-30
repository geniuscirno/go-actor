package core

import (
	"context"
	"sync"
	"time"
)

type SpawnOptions struct {
	Name string
}

type process struct {
	ctx    context.Context
	cancel context.CancelFunc

	node *node

	pid PID

	behavior ProcessBehavior

	mailbox chan Message
	exit    chan struct{}

	parent *process

	mu       sync.RWMutex
	children map[string]*process
}

func (p *process) Spawn(behavior ProcessBehavior, opts *SpawnOptions) (Process, error) {
	return p.node.spawn(p, behavior, opts)
}

func (p *process) Self() PID {
	return p.pid
}

func (p *process) Context() context.Context {
	return p.ctx
}

func (p *process) Behavior() ProcessBehavior {
	return p.behavior
}

func (p *process) IsAlive() bool {
	return p.ctx.Err() == nil
}

func (p *process) Send(to PID, message interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return p.SendCtx(ctx, to, message)
}

func (p *process) SendCtx(ctx context.Context, to PID, message interface{}) error {
	if val, ok := message.(Message); ok {
		return p.node.SendMessage(ctx, to, val)
	}
	return p.node.SendMessage(ctx, to, Message{
		From: p.Self(),
		Data: message,
	})
}

func (p *process) ProcessChannels() ProcessChannels {
	return ProcessChannels{
		Mailbox: p.mailbox,
		Exit:    p.exit,
	}
}

func (p *process) Kill() {
	p.cancel()
}

func (p *process) Stop() error {
	if !p.IsAlive() {
		return nil
	}

	select {
	case p.exit <- struct{}{}:
	default:
	}

	go func() {
		timer := time.NewTimer(time.Second * 30)
		defer timer.Stop()
		select {
		case <-p.ctx.Done():
		case <-timer.C:
			p.cancel()
		}
	}()
	return nil
}

func (p *process) Wait() {
	if p.IsAlive() {
		<-p.ctx.Done()
	}
}

func (p *process) Parent() Process {
	if p.parent == nil {
		return nil
	}
	return p.parent
}

func (p *process) addChild(child *process) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.children == nil {
		p.children = make(map[string]*process)
	}

	p.children[child.pid.ID] = child
}

func (p *process) deleteChild(child *process) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.children, child.pid.ID)
}

func (p *process) Children() []Process {
	p.mu.RLock()
	defer p.mu.RUnlock()

	children := make([]Process, 0, len(p.children))
	for _, child := range p.children {
		children = append(children, child)
	}
	return children
}

func (p *process) StopChildren() error {
	for _, child := range p.Children() {
		child.Stop()
		child.Wait()
	}
	return nil
}
