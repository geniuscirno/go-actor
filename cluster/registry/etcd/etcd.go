package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/geniuscirno/go-actor/cluster/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Option func(*options)

type options struct {
	namespace string
	ttl       time.Duration
}

func defaultOptions() options {
	return options{
		namespace: "node",
		ttl:       time.Second * 60,
	}
}

func Namespace(ns string) Option {
	return func(o *options) {
		o.namespace = ns
	}
}

func TTL(ttl time.Duration) Option {
	return func(o *options) {
		o.ttl = ttl
	}
}

type Registry struct {
	opts   options
	ctx    context.Context
	cancel context.CancelFunc

	client  *clientv3.Client
	lease   clientv3.Lease
	leaseID clientv3.LeaseID
}

func New(client *clientv3.Client, opt ...Option) registry.Registrar {
	opts := defaultOptions()
	for _, o := range opt {
		o(&opts)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Registry{client: client, ctx: ctx, cancel: cancel, opts: opts}
}

func (r *Registry) Register(ctx context.Context, node *registry.Node) error {
	key := fmt.Sprintf("%s/%s", r.opts.namespace, node.Name)
	b, err := json.Marshal(node)
	if err != nil {
		return err
	}
	r.lease = clientv3.NewLease(r.client)
	grant, err := r.lease.Grant(ctx, int64(r.opts.ttl.Seconds()))
	if err != nil {
		return err
	}
	if _, err := r.client.Put(ctx, key, string(b), clientv3.WithLease(grant.ID)); err != nil {
		return err
	}
	r.leaseID = grant.ID
	return nil
}

func (r *Registry) Deregister(ctx context.Context, node *registry.Node) error {
	defer func() {
		r.cancel()
		if r.lease != nil {
			r.lease.Close()
		}
	}()

	key := fmt.Sprintf("%s/%s", r.opts.namespace, node.Name)

	if _, err := r.client.Delete(ctx, key); err != nil {
		return err
	}
	return nil
}

func (r *Registry) KeepAlive(ctx context.Context) error {
	keepAliveChan, err := r.lease.KeepAlive(ctx, r.leaseID)
	if err != nil {
		return err
	}

	for {
		select {
		case _, ok := <-keepAliveChan:
			if !ok {
				return errors.New("keepalive chan closed")
			}
		case <-r.ctx.Done():
			return nil
		}
	}
}
