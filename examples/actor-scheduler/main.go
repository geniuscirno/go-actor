package main

import (
	"fmt"
	"go-actor/actor"
	"time"
)

type hello struct {
	Who string
}

type helloActor struct {
	cancel actor.CancelFunc
}

func (state *helloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		scheduler := actor.NewTimerScheduler(context)
		state.cancel = scheduler.SendRepeatedly(context.Self(), &hello{Who: "Roger"}, time.Second*3)
	case *actor.Stopping:
		state.cancel()
	case *actor.Stopped:
	case *hello:
		fmt.Printf("Hello %v\n", msg.Who)
	}
}

func main() {
	node := actor.NewNode("node")

	p, err := node.SpawnActor(&helloActor{})
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 30)
	p.Stop()
	time.Sleep(time.Second)
}
