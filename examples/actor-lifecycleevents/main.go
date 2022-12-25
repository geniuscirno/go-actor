package main

import (
	"fmt"
	"go-actor/actor"
	"time"
)

type hello struct {
	Who string
}

type helloActor struct{}

func (state *helloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println(context.Self(), "Started, initialize actor here")
		time.Sleep(time.Second * 1800)
	case *actor.Stopping:
		fmt.Println(context.Self(), "Stopping, actor is about shutdown")
	case *actor.Stopped:
		fmt.Println(context.Self(), "Stopped, actor and its children are stopped")
	case *hello:
		fmt.Printf("%v Hello %v\n", context.Self(), msg.Who)
	}
}

func main() {
	node := actor.NewNode("node")

	p, err := node.SpawnActor(&helloActor{}, actor.Name("parent"))
	if err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ {
		child, err := p.SpawnActor(&helloActor{}, actor.Name(fmt.Sprintf("child_%d", i)))
		if err != nil {
			panic(err)
		}
		for j := 0; j < 3; j++ {
			_, err := child.SpawnActor(&helloActor{}, actor.Name(fmt.Sprintf("child_%d_%d", i, j)))
			if err != nil {
				panic(err)
			}
		}
	}

	p.Send(p.Self(), &hello{Who: "Roger"})

	time.Sleep(time.Second * 1)
	go func() {
		time.Sleep(time.Second)
		if err := p.Stop(); err != nil {
			fmt.Println(err)
		}
	}()
	p.Wait()
	time.Sleep(time.Second)
}
