package cluster

import (
	"context"
	"errors"
	"fmt"
	"github.com/geniuscirno/go-actor/cluster/registry"
	"github.com/geniuscirno/go-actor/core"
	"github.com/geniuscirno/go-actor/remote"
	"github.com/geniuscirno/go-actor/remote/attributes"
	"net"
	"sync"
	"time"
)

type Options struct {
	registrar registry.Registrar
	address   string
	port      int
}

type Option func(opts *Options)

func Registrar(registrar registry.Registrar) Option {
	return func(opts *Options) {
		opts.registrar = registrar
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

	server *remote.Server

	mu        sync.RWMutex
	endpoints map[string]*remote.Endpoint
}

func NewCluster(node core.Node, opt ...Option) *Cluster {
	opts := Options{}
	for _, o := range opt {
		o(&opts)
	}
	cluster := &Cluster{node: node, opts: opts}
	if cluster.opts.address == "" {
		address, err := localIPV4Addr()
		if err != nil {
			panic(err)
		}
		cluster.opts.address = address
	}

	cluster.server = remote.NewServer(node, fmt.Sprintf("%s:%d", cluster.opts.address, cluster.opts.port))
	cluster.endpoints = make(map[string]*remote.Endpoint)

	if err := cluster.server.Start(); err != nil {
		panic(err)
	}
	node.Join(cluster)

	if cluster.opts.registrar != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		if err := cluster.opts.registrar.Register(ctx, &registry.Node{Name: node.Name(), Address: cluster.server.Address()}); err != nil {
			panic(err)
		}
		go cluster.opts.registrar.KeepAlive(context.TODO())
	}

	return cluster
}

func (c *Cluster) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	if c.opts.registrar != nil {
		c.opts.registrar.Deregister(ctx, &registry.Node{Name: c.node.Name(), Address: c.server.Address()})
	}
	c.server.Stop()
}

func (c *Cluster) SendMessage(ctx context.Context, to core.PID, message core.Message) error {
	ep, ok := c.getEndpoint(to.Node)
	if !ok {
		return errors.New("no node")
	}

	return ep.SendMessage(ctx, to, message)
}

func (c *Cluster) StaticRoute(nodeName string, nodeAddr string, attrs *attributes.Attributes) error {
	ep, err := remote.NewEndpoint(nodeName, nodeAddr)
	if err != nil {
		return err
	}

	c.updateEndpoint(ep)
	return nil
}

func (c *Cluster) updateEndpoint(ep *remote.Endpoint) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.endpoints[ep.NodeName] = ep
}

func (c *Cluster) UpdateEndpoints(endpoints []*remote.Endpoint) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k := range c.endpoints {
		delete(c.endpoints, k)
	}
	for _, ep := range endpoints {
		c.endpoints[ep.NodeName] = ep
	}
}

func (c *Cluster) getEndpoint(nodeName string) (*remote.Endpoint, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

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
