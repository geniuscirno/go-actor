package actor

import (
	"context"
	"time"

	"github.com/geniuscirno/go-actor/core"
)

type Process interface {
	core.Process
	SpawnActor(actor Actor, opt ...SpawnOption) (Process, error)
	CallCtx(ctx context.Context, to PID, message interface{}) *Future
	Call(to PID, message interface{}) *Future
}

type actorProcess struct {
	core.Process
}

func (p *actorProcess) SpawnActor(actor Actor, opt ...SpawnOption) (Process, error) {
	opts := &core.SpawnOptions{}
	for _, o := range opt {
		o(opts)
	}
	process, err := p.Process.Spawn(&actorBehavior{
		actor: actor,
	}, opts)
	if err != nil {
		return nil, err
	}
	return &actorProcess{Process: process}, nil
}

func (p *actorProcess) Call(to PID, message interface{}) *Future {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return p.CallCtx(ctx, to, message)
}

func (p *actorProcess) CallCtx(ctx context.Context, to PID, message interface{}) *Future {
	var timeout time.Duration
	if deadline, ok := ctx.Deadline(); ok {
		timeout = time.Until(deadline)
	}
	future, err := NewFuture(p, timeout)
	if err != nil {
		panic(err)
	}

	future.Send(to, message)
	return future
}
