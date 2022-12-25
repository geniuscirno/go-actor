package actor

import (
	"time"
)

type CancelFunc func()

type TimerScheduler struct {
	process Process
}

func NewTimerScheduler(process Process) *TimerScheduler {
	return &TimerScheduler{process: process}
}

func (s *TimerScheduler) SendOnce(to PID, message interface{}, delay time.Duration) CancelFunc {
	t := time.AfterFunc(delay, func() {
		s.process.Send(to, message)
	})
	return func() {
		t.Stop()
	}
}

func (s *TimerScheduler) SendRepeatedly(to PID, message interface{}, interval time.Duration) CancelFunc {
	t := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-t.C:
				s.process.Send(to, message)
			}
		}
	}()
	return func() {
		t.Stop()
	}
}
