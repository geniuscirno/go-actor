package core

const (
	MessageFlagResponse = 1
)

type Message struct {
	From      PID
	RequestID int64
	Flag      int32
	Data      interface{}
}

func (m *Message) TestFlag(flag int32) bool {
	return m.Flag&flag != 0
}
