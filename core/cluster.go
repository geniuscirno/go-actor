package core

import "context"

type Cluster interface {
	SendMessage(ctx context.Context, to PID, message Message) error
}
