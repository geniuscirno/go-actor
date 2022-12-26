package actor

import (
	"github.com/geniuscirno/go-actor/core"
	"sync"
	"time"
)

type Future struct {
	core.Process
	cond   *sync.Cond
	done   bool
	result interface{}
	err    error
}

func NewFuture(process Process, timeout time.Duration) (*Future, error) {
	future := &Future{cond: sync.NewCond(&sync.Mutex{})}

	fp := &futureProcess{
		Future: future,
	}

	p, err := process.Spawn(fp, &core.SpawnOptions{})
	if err != nil {
		return nil, err
	}
	fp.Future.Process = p

	if timeout > 0 {
		go func() {
			timer := time.NewTimer(timeout)
			defer timer.Stop()

			select {
			case <-timer.C:
				future.SetErr(core.ErrTimeout)
			}
		}()
	}

	return future, nil
}

func (f *Future) Wait() {
	f.cond.L.Lock()
	for !f.done {
		f.cond.Wait()
	}
	f.cond.L.Unlock()
}

func (f *Future) SetResult(result interface{}) {
	f.result = result
	f.Close()
}

func (f *Future) Result() (interface{}, error) {
	f.Wait()
	return f.result, f.err
}

func (f *Future) SetErr(err error) {
	f.err = err
	f.Close()
}

func (f *Future) Err() error {
	f.Wait()
	return f.err
}

func (f *Future) Close() {
	f.cond.L.Lock()
	if f.done {
		f.cond.L.Unlock()
		return
	}
	f.done = true

	f.cond.L.Unlock()
	f.cond.Signal()
}

type futureProcess struct {
	*Future
}

func (fp *futureProcess) ProcessLoop(process core.Process) error {
	channels := process.ProcessChannels()
	select {
	case msg := <-channels.Mailbox:
		if err, ok := msg.Data.(error); ok {
			fp.SetErr(err)
		} else {
			fp.SetResult(msg.Data)
		}
	case <-channels.Exit:
		return nil
	case <-process.Context().Done():
		return process.Context().Err()
	}
	return nil
}
