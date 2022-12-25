package main

import (
	"github.com/geniuscirno/go-actor/actor"
	"github.com/geniuscirno/go-actor/cluster"
	"github.com/geniuscirno/go-actor/examples/cluster-echo/echo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	node := actor.NewNode("node2")
	c := cluster.NewCluster(node, "localhost:9701")
	c.StaticRoute("node1", "localhost:9700", nil)

	p, err := node.SpawnActor(actor.ActorFunc(func(c actor.Context) {
		switch msg := c.Message().(type) {
		case *echo.Echo:
			log.Printf("echo request from %s: %s\n", c.From(), msg.Message)
			c.Send(c.From(), msg)
		}
	}), actor.Name("echo-server"))
	if err != nil {
		panic(err)
	}
	defer p.Stop()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT)
	<-sig
}
