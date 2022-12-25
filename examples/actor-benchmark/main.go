package main

import (
	"github.com/geniuscirno/go-actor/actor"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	N = 1000000
)

type ping struct{}
type pong struct{}

var (
	_ping = &ping{}
	_pong = &pong{}
)

func main() {
	node := actor.NewNode("node")

	wg := sync.WaitGroup{}
	p1, err := node.SpawnActor(actor.ActorFunc(func(c actor.Context) {
		switch c.Message().(type) {
		case *pong:
			wg.Done()
		}
	}))
	if err != nil {
		panic(err)
	}

	p2, err := node.SpawnActor(actor.ActorFunc(func(c actor.Context) {
		switch c.Message().(type) {
		case *ping:
			c.Send(c.From(), _pong)
		}
	}))

	wg.Add(N)
	go func() {
		ts := time.Now()
		for i := 0; i < N; i++ {
			p1.Send(p2.Self(), _ping)
		}
		wg.Wait()
		elapsed := time.Since(ts)
		qps := int32(float32(N*2) / (float32(elapsed) / float32(time.Second)))
		log.Printf("elapsed: %v\n", elapsed)
		log.Printf("qps: %v\n", qps)
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT)
	<-sig
}
