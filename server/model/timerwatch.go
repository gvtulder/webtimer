package model

import (
	"time"
)

type TimerWatch struct {
	C         chan TimerState
	timer     *Timer
	lastState TimerState
	ticker    *time.Ticker
	done      chan struct{}
}

func NewTimerWatch(timer *Timer) *TimerWatch {
	return &TimerWatch{timer: timer, C: make(chan TimerState), done: make(chan struct{})}
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
					w.ticker.Stop()
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
		w.ticker.Stop()
		w.done <- struct{}{}
	}
}
