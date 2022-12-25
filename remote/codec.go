package remote

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type MessageCodec interface {
	Marshal(m proto.Message) (*anypb.Any, error)
	Unmarshal(v *anypb.Any) (proto.Message, error)
}

type protoCodec struct{}

func (c *protoCodec) Marshal(m proto.Message) (*anypb.Any, error) {
	return anypb.New(m)
}

func (c *protoCodec) Unmarshal(v *anypb.Any) (proto.Message, error) {
	return anypb.UnmarshalNew(v, proto.UnmarshalOptions{})
}

func Marshal(m proto.Message) (*anypb.Any, error) {
	return defaultCodec.Marshal(m)
}

func Unmarshal(v *anypb.Any) (proto.Message, error) {
	return defaultCodec.Unmarshal(v)
}

var defaultCodec = &protoCodec{}
