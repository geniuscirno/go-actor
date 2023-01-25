package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/geniuscirno/go-actor/cluster/registry"
	"time"
)

type Address struct {
	Addr string
	Name string
}

type State struct {
	Addresses []Address
}

type Cluster interface {
	UpdateState(State) error
}

type Resolver interface {
	Close()
}

type resolver struct {
	ctx    context.Context
	cancel context.CancelFunc
	d      registry.Discovery
	w      registry.Watcher
	c      Cluster
}

func New(d registry.Discovery, c Cluster) (Resolver, error) {
	r := &resolver{d: d, c: c}
	r.ctx, r.cancel = context.WithCancel(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	w, err := r.d.Watch(ctx)
	if err != nil {
		return nil, err
	}
	r.w = w
	r.c = c
	go r.watch()
	return r, nil
}

func (r *resolver) watch() {
	for {
		nodes, err := r.w.Next()
		if err != nil {
			return
		}
		b, _ := json.Marshal(nodes)
		fmt.Println("resolver: watch", string(b))

		addrs := make([]Address, 0, len(nodes))
		for _, node := range nodes {
			addrs = append(addrs, Address{Name: node.Name, Addr: node.Address})
		}
		r.c.UpdateState(State{Addresses: addrs})
	}
}

func (r *resolver) Close() {
	r.w.Stop()
}
