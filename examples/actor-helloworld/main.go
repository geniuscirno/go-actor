package main

import (
	"fmt"
	"github.com/geniuscirno/go-actor/actor"
	"time"
)

type hello struct {
	Who string
}

type helloActor struct{}

func (state *helloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *hello:
		fmt.Printf("Hello %s\n", msg)
		time.Sleep(time.Second * 3)
		context.Send(context.Self(), &hello{Who: "self"})
	}
}

func main() {
	node := actor.NewNode("node")
	p, err := node.SpawnActor(&helloActor{})
	if err != nil {
		panic(err)
	}
	p.Send(p.Self(), &hello{Who: "Roger"})
	time.Sleep(time.Second * 30)
}
