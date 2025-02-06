package model

import (
	"math"
	"time"
)

type TimerState struct {
	Active    bool
	Black     bool
	Countdown bool
	Running   bool
	Remaining int64
}

type Timer struct {
	C                chan struct{}
	black            bool
	countdown        bool
	remainingAtStart int64
	startTime        int64
}

func NewTimer() *Timer {
	return &Timer{C: make(chan struct{})}
}

func (t *Timer) Active() bool {
	return t.startTime != 0 || t.remainingAtStart != 0
}

func (t *Timer) Black() bool {
	return t.black
}

func (t *Timer) Countdown() bool {
	return t.countdown
}

func (t *Timer) Running() bool {
	return t.startTime != 0
}

func (t *Timer) Remaining() int64 {
	var r int64
	if t.startTime != 0 {
		r = t.remainingAtStart - (time.Now().UnixMilli() - t.startTime)
	} else {
		r = t.remainingAtStart
	}
	return int64(math.Round(float64(r) / 1000))
}

func (t *Timer) State() TimerState {
	return TimerState{
		Active:    t.Active(),
		Black:     t.Black(),
		Countdown: t.Countdown(),
		Running:   t.Running(),
		Remaining: t.Remaining(),
	}
}

func (t *Timer) SetRemaining(s int64) {
	if t.startTime != 0 {
		t.startTime = time.Now().UnixMilli()
	}
	t.remainingAtStart = s * 1000
	if t.Remaining() > 0 {
		t.countdown = true
	}
	t.C <- struct{}{}
}

func (t *Timer) AddRemaining(s int64) {
	t.remainingAtStart += s * 1000
	if t.Remaining() > 0 {
		t.countdown = true
	}
	t.C <- struct{}{}
}

func (t *Timer) Start() {
	if t.startTime == 0 {
		t.startTime = time.Now().UnixMilli()
		t.C <- struct{}{}
	}
}

func (t *Timer) Pause() {
	if t.Running() {
		t.remainingAtStart = t.remainingAtStart - (time.Now().UnixMilli() - t.startTime)
		t.startTime = 0
		t.C <- struct{}{}
	}
}

func (t *Timer) Reset() {
	if t.Active() {
		t.startTime = 0
		t.remainingAtStart = 0
		t.countdown = false
		t.C <- struct{}{}
	}
}

func (t *Timer) ToggleBlack() {
	t.black = !t.black
	t.C <- struct{}{}
}

func (t *Timer) BlackOn() {
	if !t.black {
		t.black = true
		t.C <- struct{}{}
	}
}

func (t *Timer) BlackOff() {
	if t.black {
		t.black = false
		t.C <- struct{}{}
	}
}
