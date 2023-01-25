package actor

import (
	"github.com/geniuscirno/go-actor/core"
)

type Context interface {
	Process
	Message() interface{}
	Reply(message interface{}) error
	Error(err error) error
	HandleCall(reply interface{}, err error) error
	From() PID
}

type actorContext struct {
	message interface{}
	*actorProcess
}

func newActorContext(process *actorProcess, message interface{}) *actorContext {
	return &actorContext{actorProcess: process, message: message}
}

func (c *actorContext) Message() interface{} {
	if m, ok := c.message.(core.Message); ok {
		return m.Data
	}
	return c.message
}

func (c *actorContext) Reply(message interface{}) error {
	if m, ok := c.message.(core.Message); ok {
		return c.Send(m.From, core.Message{
			From:      c.Self(),
			RequestID: m.RequestID,
			Data:      message,
			Flag:      core.MessageFlagResponse,
		})
	}
	return nil
}

func (c *actorContext) Error(err error) error {
	return c.Reply(err)
}

func (c *actorContext) From() PID {
	if m, ok := c.message.(core.Message); ok {
		return m.From
	}
	return PID{}
}

func (c *actorContext) HandleCall(reply interface{}, err error) error {
	if err != nil {
		return c.Error(err)
	}
	return c.Send(c.From(), core.Message{From: c.Self(), Data: reply})
}
