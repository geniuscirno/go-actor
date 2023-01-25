package remote

import (
	"context"
	"errors"
	"fmt"
	"github.com/geniuscirno/go-actor/core"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
)

type Server struct {
	node core.Node

	addr string

	s *grpc.Server
}

func NewServer(node core.Node, addr string) *Server {
	return &Server{node: node, addr: addr}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	sp := strings.SplitN(s.addr, ":", 2)
	s.addr = fmt.Sprintf("%s:%d", sp[0], lis.Addr().(*net.TCPAddr).Port)

	s.s = grpc.NewServer()
	RegisterRemoteServer(s.s, s)
	go s.s.Serve(lis)
	return nil
}

func (s *Server) Address() string {
	return s.addr
}

func (s *Server) Stop() {
	s.s.Stop()
}

func (s *Server) OnMessage(ctx context.Context, in *OnMessageRequest) (*OnMessageReply, error) {
	if in.To.Node != s.node.Name() {
		return nil, errors.New("not this node")
	}

	to := core.PID{Node: in.To.Node, ID: in.To.Id}

	data, err := Unmarshal(in.Message.Data)
	if err != nil {
		return nil, err
	}
	log.Printf("cluster: recv message from %s@%s to %v: %v\n", in.Message.From.Id, in.Message.From.Node, to, data)

	if err := s.node.SendMessage(ctx, to, core.Message{
		From:      core.PID{Node: in.Message.From.Node, ID: in.Message.From.Id},
		RequestID: in.Message.RequestId,
		Data:      data,
		Flag:      in.Message.Flag,
	}); err != nil {
		return nil, err
	}

	return &OnMessageReply{}, nil
}

func (s *Server) mustEmbedUnimplementedRemoteServer() {}
