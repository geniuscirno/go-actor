package core

import (
	"context"
	"fmt"
	"strings"
)

type PID struct {
	Node string
	ID   string
}

func PIDFromString(s string) (PID, error) {
	sp := strings.Split(s, "@")
	if len(sp) != 2 {
		return PID{}, fmt.Errorf("invalid pid string: %s", s)
	}
	return PID{Node: sp[0], ID: sp[1]}, nil
}

func (p PID) String() string {
	return fmt.Sprintf("%s@%s", p.ID, p.Node)
}

type Process interface {
	Context() context.Context
	Self() PID
	IsAlive() bool
	Send(to PID, message interface{}) error
	SendCtx(ctx context.Context, to PID, message interface{}) error
	Spawn(behavior ProcessBehavior, opts *SpawnOptions) (Process, error)
	Behavior() ProcessBehavior
	ProcessChannels() ProcessChannels
	Stop() error
	Wait()
	Kill()

	Parent() Process
	Children() []Process
	StopChildren() error
}

type ProcessBehavior interface {
	ProcessLoop(Process) error
}

type ProcessChannels struct {
	Mailbox <-chan Message
	Exit    <-chan struct{}
}
