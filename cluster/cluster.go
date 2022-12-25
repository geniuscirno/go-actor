package cluster

import (
	"context"
	"errors"
	"go-actor/core"
	"go-actor/remote"
	"go-actor/remote/attributes"
	"sync"
)

type Cluster struct {
	node core.Node

	server *remote.Server

	mu        sync.RWMutex
	endpoints map[string]*remote.Endpoint
}

func NewCluster(node core.Node, addr string) *Cluster {
	cluster := &Cluster{node: node}
	cluster.server = remote.NewServer(node, addr)
	cluster.endpoints = make(map[string]*remote.Endpoint)

	if err := cluster.server.Start(); err != nil {
		panic(err)
	}
	node.Join(cluster)
	return cluster
}

func (c *Cluster) Stop() {
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
