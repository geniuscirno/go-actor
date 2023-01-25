package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/geniuscirno/go-actor/cluster/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type watcher struct {
	key    string
	ctx    context.Context
	cancel context.CancelFunc

	nodes map[string]*registry.Node

	watchChan clientv3.WatchChan
	watcher   clientv3.Watcher
	kv        clientv3.KV
	first     bool
}

func newWatcher(ctx context.Context, key string, client *clientv3.Client) (*watcher, error) {
	w := &watcher{
		key:     key,
		watcher: clientv3.NewWatcher(client),
		kv:      clientv3.NewKV(client),
		nodes:   make(map[string]*registry.Node),
		first:   true,
	}
	w.ctx, w.cancel = context.WithCancel(ctx)

	resp, err := w.kv.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, kv := range resp.Kvs {
		n := &registry.Node{}
		if err := json.Unmarshal(kv.Value, n); err != nil {
			continue
		}
		w.nodes[string(kv.Key)] = n
	}

	w.watchChan = w.watcher.Watch(w.ctx, key, clientv3.WithPrefix(), clientv3.WithRev(resp.Header.Revision+1))
	return w, nil
}

func (w *watcher) Next() ([]*registry.Node, error) {
	if w.first {
		w.first = false
		return nodeMapToSlice(w.nodes), nil
	}
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case wch, ok := <-w.watchChan:
		if !ok {
			return nil, errors.New("watch chan closed")
		}

		for _, event := range wch.Events {
			n := &registry.Node{}
			switch event.Type {
			case clientv3.EventTypePut:
				if err := json.Unmarshal(event.Kv.Value, n); err != nil {
					continue
				}
				w.nodes[string(event.Kv.Key)] = n
			case clientv3.EventTypeDelete:
				delete(w.nodes, string(event.Kv.Key))
			}
		}

		return nodeMapToSlice(w.nodes), nil
	}
}

func nodeMapToSlice(nodes map[string]*registry.Node) []*registry.Node {
	results := make([]*registry.Node, 0, len(nodes))
	for _, node := range nodes {
		results = append(results, node)
	}
	return results
}

func (w *watcher) Stop() error {
	w.cancel()
	return w.watcher.Close()
}
