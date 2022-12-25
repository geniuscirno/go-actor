package actor

import (
	"fmt"
	"go-actor/actor/future"
	"go-actor/core"
	"sync"
)

type PID = core.PID

type Actor interface {
	Receive(c Context)
}

//goland:noinspection GoNameStartsWithPackageName
type ActorFunc func(c Context)

func (f ActorFunc) Receive(c Context) {
	f(c)
}

type actorBehavior struct {
	mu        sync.RWMutex
	futures   map[int64]*future.Future
	requestId int64

	actor Actor
}

func (b *actorBehavior) ProcessLoop(process core.Process) error {
	actorProcess := &actorProcess{Process: process}

	b.handleStarted(actorProcess)
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
		b.handleTerminate(actorProcess)
	}()

	channels := process.ProcessChannels()
	for {
		select {
		case <-process.Context().Done():
			return process.Context().Err()
		case <-channels.Exit:
			b.handleStop(actorProcess)
			return nil
		case message := <-channels.Mailbox:
			if message.TestFlag(core.MessageFlagResponse) {
				b.handleReply(message)
				continue
			}

			b.actor.Receive(newActorContext(actorProcess, message))
		}
	}
}

func (b *actorBehavior) handleStarted(process *actorProcess) {
	b.actor.Receive(newActorContext(process, startedMessage))
}

func (b *actorBehavior) handleStop(process *actorProcess) {
	b.actor.Receive(newActorContext(process, stoppingMessage))
	process.StopChildren()
}

func (b *actorBehavior) handleTerminate(process *actorProcess) {
	b.actor.Receive(newActorContext(process, stoppedMessage))
}

func (b *actorBehavior) handleReply(message core.Message) {
	if message.RequestID == 0 {
		return
	}

	b.mu.Lock()
	f, ok := b.futures[message.RequestID]
	if !ok {
		b.mu.Unlock()
		return
	}
	delete(b.futures, message.RequestID)
	b.mu.Unlock()

	if err, ok := message.Data.(error); ok {
		f.SetErr(err)
		return
	}
	f.SetResult(message.Data)
}
