package main

import (
	"bufio"
	"fmt"
	"go-actor/actor"
	"go-actor/cluster"
	"go-actor/examples/cluster-chat/chat"
	"log"
	"os"
)

func main() {
	node := actor.NewNode("client")
	c := cluster.NewCluster(node, "localhost:9701")
	c.StaticRoute("server", "localhost:9700", nil)

	server := actor.PID{Node: "server", ID: "chat-server"}

	p, err := node.SpawnActor(actor.ActorFunc(func(c actor.Context) {
		switch msg := c.Message().(type) {
		case *chat.Connected:
			log.Println(msg.Message)
		case *chat.SayReply:
			log.Printf("%s: %s\n", msg.Username, msg.Message)
		}
	}))
	if err != nil {
		panic(err)
	}
	defer p.Stop()

	username := fmt.Sprintf("%s@%s", p.Self().ID, p.Self().Node)
	p.Send(server, &chat.Connect{})

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		p.Send(server, &chat.SayRequest{
			Username: username,
			Message:  text,
		})
	}
}
