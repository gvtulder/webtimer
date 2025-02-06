// Copyright (C) 2025 Gijs van Tulder
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package server

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
