package main

import (
	"github.com/geniuscirno/go-actor/actor"
	"github.com/geniuscirno/go-actor/cluster"
	"github.com/geniuscirno/go-actor/examples/cluster-benchmark/messages"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	N = 100000
)

func main() {
	node := actor.NewNode("node1")
	c := cluster.NewCluster(node, "localhost:9700")
	c.StaticRoute("node2", "localhost:9701", nil)

	wg := &sync.WaitGroup{}

	p, err := node.SpawnActor(actor.ActorFunc(func(c actor.Context) {
		switch c.Message().(type) {
		case *messages.Pong:
			wg.Done()
		}
	}))
	if err != nil {
		panic(err)
	}
	defer p.Stop()

	to := actor.PID{Node: "node2", ID: "server"}
	wg.Add(N)
	go func() {
		log.Println("start to send")
		ts := time.Now()
		ping := &messages.Ping{}
		for i := 0; i < N; i++ {
			if err := p.Send(to, ping); err != nil {
				panic(err)
			}
		}
		wg.Wait()
		elapsed := time.Since(ts)
		qps := float32(N*2) / (float32(elapsed) / float32(time.Second))
		log.Printf("elapsed: %v\n", elapsed)
		log.Printf("qps: %v\n", qps)
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT)
	<-sig
}
