package model

import (
	"time"
)

type TimerWatch struct {
	C         chan TimerState
	timer     *Timer
	lastState TimerState
	ticker    *time.Ticker
	done      chan bool
}

func NewTimerWatch(timer *Timer) *TimerWatch {
	return &TimerWatch{timer: timer, C: make(chan TimerState), done: make(chan bool)}
}

func (w *TimerWatch) Start() {
	if w.ticker == nil {
		w.ticker = time.NewTicker(200 * time.Millisecond)
		w.HandleUpdate(w.timer.State())
		go func() {
			for {
				select {
				case <-w.ticker.C:
					w.HandleUpdate(w.timer.State())
				case <-w.timer.C:
					w.HandleUpdate(w.timer.State())
				case <-w.done:
					return
				}
			}
		}()
	}
}

func (w *TimerWatch) HandleUpdate(s TimerState) {
	if w.lastState != s {
		w.lastState = s
		w.C <- s
	}
}

func (w *TimerWatch) Stop() {
	if w.ticker != nil {
		w.done <- false
		w.ticker.Stop()
	}
}
