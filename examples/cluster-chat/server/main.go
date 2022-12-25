package main

import (
	"go-actor/actor"
	"go-actor/cluster"
	"go-actor/examples/cluster-chat/chat"
	"log"
)

func notifyAll(context actor.Context, clients map[actor.PID]struct{}, message interface{}) {
	for client := range clients {
		context.Send(client, message)
	}
}

func main() {
	node := actor.NewNode("server")
	c := cluster.NewCluster(node, "localhost:9700")
	c.StaticRoute("client", "localhost:9701", nil)

	clients := make(map[actor.PID]struct{})

	p, err := node.SpawnActor(actor.ActorFunc(func(c actor.Context) {
		switch msg := c.Message().(type) {
		case *chat.Connect:
			log.Printf("Client %v connected", c.From())
			clients[c.From()] = struct{}{}
			c.Send(c.From(), &chat.Connected{Message: "Welcome!"})
		case *chat.SayRequest:
			notifyAll(c, clients, &chat.SayReply{
				Username: msg.Username,
				Message:  msg.Message,
			})
		}
	}), actor.Name("chat-server"))
	if err != nil {
		panic(err)
	}
	defer p.Stop()

	select {}
}
