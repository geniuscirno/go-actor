package core

import "context"

type Cluster interface {
	Spawn(behavior ProcessBehavior, opts *SpawnOptions) (Process, error)
	SendMessage(ctx context.Context, to PID, message Message) error
	Stop()
}
