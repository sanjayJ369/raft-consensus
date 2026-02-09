package simple

import (
	"sync"
	"time"
)

type Timer struct {
	sync.RWMutex
	duration time.Duration
	start    time.Time
	timer    *time.Timer
	callback func()
	active   bool
}

func NewSimpleTimer() *Timer {
	return &Timer{}
}

func (t *Timer) Start(d time.Duration, f func()) {
	t.Lock()
	defer t.Unlock()

	// stop existing timer if running
	if t.timer != nil {
		t.timer.Stop()
	}

	t.duration = d
	t.callback = f
	t.start = time.Now()
	t.active = true

	t.timer = time.AfterFunc(d, func() {
		t.Lock()
		t.active = false
		t.Unlock()
		f()
	})
}

func (t *Timer) Stop() {
	t.Lock()
	defer t.Unlock()

	if t.timer != nil {
		t.timer.Stop()
	}
	t.active = false
}

func (t *Timer) Restart() {
	t.Lock()
	defer t.Unlock()

	if t.callback == nil {
		return
	}

	if t.timer != nil {
		t.timer.Stop()
	}

	t.start = time.Now()
	t.active = true

	t.timer = time.AfterFunc(t.duration, func() {
		t.Lock()
		t.active = false
		t.Unlock()
		t.callback()
	})
}

func (t *Timer) Elapsed() time.Duration {
	if !t.active {
		return 0
	}

	t.RLock()
	defer t.RUnlock()

	if !t.active {
		return 0
	}
	return time.Since(t.start)
}

func (t *Timer) Duration() time.Duration {
	t.RLock()
	defer t.RUnlock()
	return t.duration
}
