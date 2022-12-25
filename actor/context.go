package actor

import (
	"go-actor/core"
)

type Context interface {
	Process
	Message() interface{}
	Reply(message interface{}) error
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

func (c *actorContext) From() PID {
	if m, ok := c.message.(core.Message); ok {
		return m.From
	}
	return PID{}
}
