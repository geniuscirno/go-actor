package main

import (
	"fmt"
	"github.com/geniuscirno/go-actor/actor"
	"github.com/geniuscirno/go-actor/cluster"
	"github.com/geniuscirno/go-actor/examples/cluster-echo/echo"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	node := actor.NewNode("node1")
	c := cluster.NewCluster(node, "localhost:9700")
	c.StaticRoute("node2", "localhost:9701", nil)

	p, err := node.SpawnActor(actor.ActorFunc(func(c actor.Context) {
		switch msg := c.Message().(type) {
		case *echo.Echo:
			fmt.Printf("echo reply from %s: %s\n", c.From(), msg.Message)
		}
	}))
	if err != nil {
		panic(err)
	}
	defer p.Stop()

	to := actor.PID{Node: "node2", ID: "echo-server"}
	if err := p.Send(to, &echo.Echo{Message: "hello"}); err != nil {
		panic(err)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT)
	<-sig
}
