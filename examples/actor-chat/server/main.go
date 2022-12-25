package main

import (
	"fmt"
	"github.com/geniuscirno/go-actor/actor"
	"github.com/geniuscirno/go-actor/examples/actor-chat/chat"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	clients = make(map[*Client]struct{})
)

type Client struct {
	Username string
	conn     *websocket.Conn
	actor.Process
}

func serveWs(serverProcess actor.Process, w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	clientProcess, err := serverProcess.SpawnActor(actor.ActorFunc(func(c actor.Context) {
		switch msg := c.Message().(type) {
		case *chat.SayRequest:
			c.Send(serverProcess.Self(), msg)
		case *chat.SayReply:
			b, err := proto.Marshal(msg)
			if err != nil {
				panic(err)
			}
			conn.WriteMessage(websocket.BinaryMessage, b)
		case *actor.Stopped:
			log.Println(c.Self(), "Stopped")
		}
	}))
	if err != nil {
		panic(err)
	}

	client := &Client{conn: conn, Username: username, Process: clientProcess}
	clientProcess.Send(serverProcess.Self(), &clientConnected{client: client})

	go func() {
		defer func() {
			conn.Close()
			clientProcess.Send(serverProcess.Self(), &clientDisconnect{client: client})
			clientProcess.Stop()
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("error: ", err)
				}
				break
			}

			sayRequest := &chat.SayRequest{}
			if err := proto.Unmarshal(message, sayRequest); err != nil {
				panic(err)
			}

			clientProcess.Send(clientProcess.Self(), sayRequest)
		}
	}()
}

func notifyAll(context actor.Context, message interface{}) {
	for client := range clients {
		context.Send(client.Self(), message)
	}
}

type clientConnected struct {
	client *Client
}

type clientDisconnect struct {
	client *Client
}

func main() {
	node := actor.NewNode("server")

	p, err := node.SpawnActor(actor.ActorFunc(func(c actor.Context) {
		switch msg := c.Message().(type) {
		case *clientConnected:
			log.Println("client connected", msg.client.Username)
			clients[msg.client] = struct{}{}
			notifyAll(c, &chat.SayReply{Username: "server", Message: fmt.Sprintf("Welcome, %s!", msg.client.Username)})
		case *clientDisconnect:
			log.Println("client disconnected", msg.client.Username)
			delete(clients, msg.client)
			notifyAll(c, &chat.SayReply{Username: "server", Message: fmt.Sprintf("%s Leaved!", msg.client.Username)})
		case *chat.SayRequest:
			notifyAll(c, &chat.SayReply{
				Username: msg.Username,
				Message:  msg.Message,
			})
		case *actor.Stopped:
			log.Println(c.Self(), "Stopped")
		}
	}))
	if err != nil {
		panic(err)
	}
	defer p.Stop()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(p, w, r)
	})
	go func() {
		if err := http.ListenAndServe("localhost:9700", nil); err != nil {
			log.Println(err)
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	select {
	case <-sig:
		log.Println("exit...")
	}
}
