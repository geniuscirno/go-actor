package main

import (
	"github.com/geniuscirno/go-actor/actor"
	"github.com/geniuscirno/go-actor/cluster"
	"github.com/geniuscirno/go-actor/examples/cluster-benchmark/messages"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	node := actor.NewNode("node2")
	c := cluster.NewCluster(node, "localhost:9701")
	c.StaticRoute("node1", "localhost:9700", nil)

	pong := &messages.Pong{}
	p, err := node.SpawnActor(actor.ActorFunc(func(c actor.Context) {
		switch c.Message().(type) {
		case *messages.Ping:
			c.Send(c.From(), pong)
		}
	}), actor.Name("server"))
	if err != nil {
		panic(err)
	}
	defer p.Stop()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT)
	<-sig
}
