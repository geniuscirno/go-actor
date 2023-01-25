package actor

import (
	"fmt"
	"github.com/geniuscirno/go-actor/core"
	"strings"
)

type PID = core.PID

var ZeroPID PID

func PIDFromString(s string) (PID, error) {
	sp := strings.Split(s, "@")
	if len(sp) != 2 {
		return PID{}, fmt.Errorf("invalid pid string: %s", s)
	}
	return PID{Node: sp[0], ID: sp[1]}, nil
}

type Actor interface {
	Receive(c Context)
}

//goland:noinspection GoNameStartsWithPackageName
type ActorFunc func(c Context)

func (f ActorFunc) Receive(c Context) {
	f(c)
}

type actorBehavior struct {
	actor Actor
}

func (b *actorBehavior) ProcessLoop(process core.Process) error {
	actorProcess := &actorProcess{Process: process}
	defer func() {
		//if e := recover(); e != nil {
		//	fmt.Println(e)
		//}
		b.handleTerminate(actorProcess)
	}()

	b.handleStarted(actorProcess)

	channels := process.ProcessChannels()
	for {
		select {
		case <-process.Context().Done():
			return process.Context().Err()
		case <-channels.Exit:
			b.handleStop(actorProcess)
			return nil
		case message := <-channels.Mailbox:
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

//func (b *actorBehavior) handleReply(message core.Message) {
//	if message.RequestID == 0 {
//		return
//	}
//
//	b.mu.Lock()
//	f, ok := b.futures[message.RequestID]
//	if !ok {
//		b.mu.Unlock()
//		return
//	}
//	delete(b.futures, message.RequestID)
//	b.mu.Unlock()
//
//	if err, ok := message.Data.(error); ok {
//		f.SetErr(err)
//		return
//	}
//	f.SetResult(message.Data)
//}
