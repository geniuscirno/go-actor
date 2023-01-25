package remote

import (
	"context"
	"errors"
	"fmt"
	"github.com/geniuscirno/go-actor/core"
	"github.com/geniuscirno/go-actor/remote/attributes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

type Endpoint struct {
	Name       string
	Addr       string
	Attributes *attributes.Attributes
	conn       *grpc.ClientConn
}

func NewEndpoint(nodeName string, addr string) (*Endpoint, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Endpoint{
		Name: nodeName,
		Addr: addr,
		conn: conn,
	}, nil
}

func (ep *Endpoint) Close() error {
	return ep.conn.Close()
}

func (ep *Endpoint) SendMessage(ctx context.Context, to core.PID, message core.Message) error {
	if to.Node != ep.Name {
		return errors.New("node unmatch")
	}

	protoMessage, ok := message.Data.(proto.Message)
	if !ok {
		return fmt.Errorf("remote message must be a proto message")
	}

	data, err := Marshal(protoMessage)
	if err != nil {
		return err
	}

	client := NewRemoteClient(ep.conn)

	_, err = client.OnMessage(ctx, &OnMessageRequest{
		To: &PID{Node: to.Node, Id: to.ID},
		Message: &Message{
			From:      &PID{Node: message.From.Node, Id: message.From.ID},
			RequestId: message.RequestID,
			Data:      data,
			Flag:      message.Flag,
		},
	})
	return err
}
