// Package simpletimer is an implementation of the Timer interface.
package simpletimer

import "time"

type Timer struct {
	duration time.Duration
	stop     chan bool
	timer    *time.Timer
}

func NewSimpleTimer() *Timer {
	return &Timer{
		stop: make(chan bool),
	}
}

func (t *Timer) Stop() {
	t.timer.Stop()
}

func (t *Timer) Restart() {
	t.Reset()
}

func (t *Timer) Reset() {
	t.timer.Reset(t.duration)
}

func (t *Timer) Start(duration time.Duration, f func()) {
	t.duration = duration
	t.timer = time.NewTimer(duration)
	go func() {
		for {
			<-t.timer.C
			f()
			t.Reset()
		}
	}()
}
