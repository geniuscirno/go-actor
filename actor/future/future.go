package future

import "sync"

type Future struct {
	cond   *sync.Cond
	done   bool
	result interface{}
	err    error
}

func New() *Future {
	return &Future{cond: sync.NewCond(&sync.Mutex{})}
}

func (f *Future) Done() bool {
	f.cond.L.Lock()
	defer f.cond.L.Unlock()

	return f.done
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

func (f *Future) SetErr(err error) {
	f.err = err
	f.Close()
}

func (f *Future) Result() (interface{}, error) {
	f.Wait()

	return f.result, f.err
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
