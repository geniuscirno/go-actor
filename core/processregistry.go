package core

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type ProcessRegistry struct {
	nextPID int64

	mu        sync.RWMutex
	processes map[string]*process
}

func newProcessRegistry() *ProcessRegistry {
	return &ProcessRegistry{
		processes: make(map[string]*process),
	}
}

func (r *ProcessRegistry) NextId() string {
	next := atomic.AddInt64(&r.nextPID, 1)
	return fmt.Sprintf("$%d", next)
}

func (r *ProcessRegistry) Add(p *process) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.processes[p.pid.ID]
	if ok {
		return ErrDupProcessName
	}

	r.processes[p.pid.ID] = p
	return nil
}

func (r *ProcessRegistry) Delete(p *process) {
	r.mu.Lock()
	delete(r.processes, p.pid.ID)
	r.mu.Unlock()
}

func (r *ProcessRegistry) Get(pid PID) (*process, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.processes[pid.ID]
	if !ok {
		return nil, ErrProcessNotFound
	}
	return p, nil
}
