package main

import (
	"fmt"
	"github.com/geniuscirno/go-actor/actor"
	"time"
)

type hello struct {
	Who string
}

func Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *hello:
		time.Sleep(time.Second * 3)
		context.Reply("Hello " + msg.Who)
	}
}

func main() {
	node := actor.NewNode("node")

	p, err := node.SpawnActor(actor.ActorFunc(Receive))
	if err != nil {
		panic(err)
	}

	result, err := p.Call(p.Self(), &hello{Who: "Roger"}).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	time.Sleep(time.Second * 30)
}
