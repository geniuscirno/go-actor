package main

import (
	"bufio"
	"flag"
	"github.com/geniuscirno/go-actor/examples/actor-chat/chat"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"os"
)

var (
	username string
)

func main() {
	flag.StringVar(&username, "username", "", "")
	flag.Parse()

	if username == "" {
		_uuid, _ := uuid.NewRandom()
		username = _uuid.String()
	}

	header := http.Header{}
	header.Set("username", username)
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:9700/ws", header)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Printf("ReadMessage failed: %v\n", err)
				return
			}

			sayReply := &chat.SayReply{}
			if err := proto.Unmarshal(message, sayReply); err != nil {
				panic(err)
			}

			log.Printf("[%s]: %s\n", sayReply.Username, sayReply.Message)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		sayRequest := &chat.SayRequest{
			Username: username,
			Message:  scanner.Text(),
		}
		b, err := proto.Marshal(sayRequest)
		if err != nil {
			panic(err)
		}

		c.WriteMessage(websocket.BinaryMessage, b)
	}
}
