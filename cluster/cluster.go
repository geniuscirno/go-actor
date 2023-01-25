package cluster

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/geniuscirno/go-actor/cluster/registry"
	"github.com/geniuscirno/go-actor/cluster/resolver"
	"github.com/geniuscirno/go-actor/core"
	"github.com/geniuscirno/go-actor/remote"
	"github.com/geniuscirno/go-actor/remote/attributes"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"net"
	"sync"
	"time"
)

type Options struct {
	registrar registry.Registrar
	discovery registry.Discovery
	address   string
	port      int
}

type Option func(opts *Options)

func Registrar(registrar registry.Registrar) Option {
	return func(opts *Options) {
		opts.registrar = registrar
	}
}

func Discovery(discovery registry.Discovery) Option {
	return func(opts *Options) {
		opts.discovery = discovery
	}
}

func Address(addr string) Option {
	return func(opts *Options) {
		opts.address = addr
	}
}

func localIPV4Addr() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

type Cluster struct {
	opts Options
	node core.Node

	server   *remote.Server
	resolver resolver.Resolver

	lease   clientv3.Lease
	watcher clientv3.Watcher
	client  *clientv3.Client

	mu         sync.RWMutex
	endpoints  map[string]*remote.Endpoint
	globalPids map[string]core.PID
}

func NewCluster(node core.Node, client *clientv3.Client, opt ...Option) (*Cluster, error) {
	opts := Options{}
	for _, o := range opt {
		o(&opts)
	}
	cluster := &Cluster{node: node, opts: opts}
	if cluster.opts.address == "" {
		address, err := localIPV4Addr()
		if err != nil {
			return nil, err
		}
		cluster.opts.address = address
	}

	cluster.server = remote.NewServer(node, fmt.Sprintf("%s:%d", cluster.opts.address, cluster.opts.port))
	cluster.endpoints = make(map[string]*remote.Endpoint)
	cluster.globalPids = make(map[string]core.PID)
	cluster.client = client
	cluster.watcher = clientv3.NewWatcher(client)
	cluster.lease = clientv3.NewLease(client)

	if err := cluster.server.Start(); err != nil {
		return nil, err
	}
	node.Join(cluster)

	if cluster.opts.registrar != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		if err := cluster.opts.registrar.Register(ctx, &registry.Node{Name: node.Name(), Address: cluster.server.Address()}); err != nil {
			return nil, err
		}
		go cluster.opts.registrar.KeepAlive(context.TODO())
	}

	if cluster.opts.discovery != nil {
		r, err := resolver.New(cluster.opts.discovery, cluster)
		if err != nil {
			return nil, err
		}
		cluster.resolver = r
	}
	go cluster.watchGlobalPids()
	return cluster, nil
}

func (c *Cluster) Spawn(behavior core.ProcessBehavior, opts *core.SpawnOptions) (core.Process, error) {
	if opts.Name == "" {
		return nil, errors.New("cluster process must has a name")
	}

	s, err := concurrency.NewSession(c.client)
	if err != nil {
		return nil, err
	}
	defer s.Close()
	mu := concurrency.NewMutex(s, fmt.Sprintf("lock/%s", opts.Name))
	if err := mu.Lock(context.TODO()); err != nil {
		return nil, err
	}
	defer mu.Unlock(context.TODO())

	resp, err := c.client.Get(context.TODO(), fmt.Sprintf("globalPids/%s", opts.Name))
	if err != nil {
		panic(err)
	}
	if resp.Count > 0 {
		return nil, errors.New("already spawn on other nodes")
	}

	grant, err := c.lease.Grant(context.Background(), int64((time.Second * 15).Seconds()))
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(core.PID{
		Node: c.node.Name(),
		ID:   opts.Name,
	})
	if err != nil {
		return nil, err
	}

	if _, err := c.client.Put(context.Background(), fmt.Sprintf("globalPids/%s", opts.Name), string(b), clientv3.WithLease(grant.ID)); err != nil {
		panic(err)
	}
	go func() {
		kac, err := c.lease.KeepAlive(context.TODO(), grant.ID)
		if err != nil {
			panic(err)
		}

		for {
			select {
			case <-kac:
			}
		}
	}()

	return c.node.Spawn(behavior, opts)
}

func (c *Cluster) watchGlobalPids() {
	resp, err := c.client.Get(context.TODO(), "globalPids/", clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	for _, kv := range resp.Kvs {
		pid := core.PID{}
		if err := json.Unmarshal(kv.Value, &pid); err != nil {
			continue
		}
		log.Println("cluster: add global pid", pid)
		c.mu.Lock()
		c.globalPids[pid.ID] = pid
		c.mu.Unlock()
	}

	wch := c.watcher.Watch(context.TODO(), "globalPids/", clientv3.WithPrefix(), clientv3.WithRev(resp.Header.Revision+1))
	for {
		select {
		case ch, ok := <-wch:
			if !ok {
				return
			}

			c.mu.Lock()
			for _, event := range ch.Events {
				pid := core.PID{}
				switch event.Type {
				case clientv3.EventTypePut:
					if err := json.Unmarshal(event.Kv.Value, &pid); err != nil {
						continue
					}
					log.Println("cluster: add global pid", pid)
					c.globalPids[pid.ID] = pid
				case clientv3.EventTypeDelete:
					log.Println("cluster: delete global pid", string(event.Kv.Key))
					delete(c.globalPids, string(event.Kv.Key))
				}
			}
			c.mu.Unlock()
		}
	}
}

func (c *Cluster) Stop() {
	c.lease.Close()
	for _, pid := range c.globalPids {
		if pid.Node == c.node.Name() {
			c.client.Delete(context.TODO(), fmt.Sprintf("globalPids/%s", pid.ID))
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	if c.opts.registrar != nil {
		c.opts.registrar.Deregister(ctx, &registry.Node{Name: c.node.Name(), Address: c.server.Address()})
	}
	if c.resolver != nil {
		c.resolver.Close()
	}
	c.server.Stop()
}

func (c *Cluster) SendMessage(ctx context.Context, to core.PID, message core.Message) error {
	c.mu.RLock()
	if to.Node == "" {
		pid, ok := c.globalPids[to.ID]
		if ok {
			c.mu.RUnlock()
			to.Node = pid.Node
			return c.node.SendMessage(ctx, to, message)
		}
	}
	ep, ok := c.getEndpoint(to.Node)
	if !ok {
		c.mu.RUnlock()
		return errors.New("no node")
	}
	c.mu.RUnlock()

	return ep.SendMessage(ctx, to, message)
}

func (c *Cluster) StaticRoute(nodeName string, nodeAddr string, attrs *attributes.Attributes) error {
	ep, err := remote.NewEndpoint(nodeName, nodeAddr)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.updateEndpoint(ep)
	return nil
}

func (c *Cluster) updateEndpoint(ep *remote.Endpoint) {
	c.endpoints[ep.Name] = ep
}

func (c *Cluster) UpdateState(state resolver.State) error {
	log.Println("cluster: UpdateState", state.Addresses)
	c.mu.Lock()
	defer c.mu.Unlock()

	addrsSet := make(map[string]struct{})
	for _, a := range state.Addresses {
		if a.Name == c.node.Name() {
			continue
		}

		addrsSet[a.Name] = struct{}{}
		ep, ok := c.getEndpoint(a.Name)
		if !ok {
			// 新节点地址 (在endpoints里不存在)
			ep, err := remote.NewEndpoint(a.Name, a.Addr)
			if err != nil {
				continue
			}
			c.updateEndpoint(ep)
			continue
		}

		// 节点已经存在，但是地址已经改变了
		if ep.Addr != a.Addr {
			ep.Close()
			ep, err := remote.NewEndpoint(a.Name, a.Addr)
			if err != nil {
				continue
			}
			c.updateEndpoint(ep)
		}
	}

	for name, ep := range c.endpoints {
		if _, ok := addrsSet[name]; !ok {
			ep.Close()
			delete(c.endpoints, name)
		}
	}
	return nil
}

func (c *Cluster) getEndpointByAddr(addr string) (*remote.Endpoint, bool) {
	for _, v := range c.endpoints {
		if v.Addr == addr {
			return v, true
		}
	}
	return nil, false
}

func (c *Cluster) getEndpoint(nodeName string) (*remote.Endpoint, bool) {
	ep, ok := c.endpoints[nodeName]
	if !ok {
		return nil, false
	}
	return ep, true
}

func (c *Cluster) Endpoints() []*remote.Endpoint {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var results []*remote.Endpoint
	for _, ep := range c.endpoints {
		results = append(results, ep)
	}
	return results
}

//func (c *Cluster) FindEndpoints(attrs *attributes.Attributes) []*Endpoint {
//	c.mu.RLock()
//	defer c.mu.RUnlock()
//
//	var results []*Endpoint
//	for _, ep := range c.endpoints {
//		if ep.Attributes.Match(attrs) {
//			results = append(results, ep)
//		}
//	}
//	return results
//}
//
//func (c *Cluster) FindRandomEndpoint(attrs *attributes.Attributes) (*Endpoint, bool) {
//	c.mu.RLock()
//	defer c.mu.RUnlock()
//
//	for _, ep := range c.endpoints {
//		if ep.Attributes.Match(attrs) {
//			return ep, true
//		}
//	}
//	return nil, false
//}
