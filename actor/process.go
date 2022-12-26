package actor

import (
	"context"
	"github.com/geniuscirno/go-actor/core"
	"time"
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
		//futures: make(map[int64]*future.Future),
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
	//f := future.New()
	//behavior := p.Behavior().(*actorBehavior)
	//requestId := atomic.AddInt64(&behavior.requestId, 1)
	//behavior.mu.Lock()
	//behavior.futures[requestId] = f
	//behavior.mu.Unlock()
	//
	//if err := p.Process.SendCtx(ctx, to, core.Message{
	//	From:      p.Self(),
	//	RequestID: requestId,
	//	Data:      message,
	//}); err != nil {
	//	f.SetErr(err)
	//	return f
	//}
	//
	//if deadline, ok := ctx.Deadline(); ok {
	//	go func() {
	//		timer := time.NewTimer(time.Until(deadline))
	//		defer timer.Stop()
	//
	//		select {
	//		case <-timer.C:
	//			if !f.Done() {
	//				behavior.mu.Lock()
	//				delete(behavior.futures, requestId)
	//				behavior.mu.Unlock()
	//				f.SetErr(core.ErrTimeout)
	//			}
	//		}
	//	}()
	//}
	//if deadline, ok := ctx.Deadline(); ok {
	//	go func() {
	//		timer := time.NewTimer(time.Until(deadline))
	//		defer timer.Stop()
	//
	//		select {
	//		case <-timer.C:
	//			if !f
	//
	//		}
	//	}()
	//
	//}
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
