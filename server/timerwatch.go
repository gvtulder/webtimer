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

// A TimerWatch follows a Timer and sends regular updates when the time or state changes.
//
// State changes are sent to channel C.
// For a running timer, the TimerWatch will send a timer updates at least once a second.
type TimerWatch struct {
	C         chan TimerState // channel for timer updates
	timer     *Timer          // the Timer to be observed
	lastState TimerState      // the last state of the timer
	ticker    *time.Ticker    // ti
	done      chan struct{}   // channel to stop the TimerWatch
}

// NewTimerWatch initializes a new TimerWatch for the given timer.
func NewTimerWatch(timer *Timer) *TimerWatch {
	return &TimerWatch{timer: timer, C: make(chan TimerState), done: make(chan struct{})}
}

// Start starts a goroutine to observe the timer and send regular updates.
func (w *TimerWatch) Start() {
	if w.ticker == nil {
		w.ticker = time.NewTicker(200 * time.Millisecond)
		w.handleUpdate(w.timer.State())
		go func() {
			for {
				select {
				case <-w.ticker.C:
					w.handleUpdate(w.timer.State())
				case <-w.timer.C:
					w.handleUpdate(w.timer.State())
				case <-w.done:
					w.ticker.Stop()
					return
				}
			}
		}()
	}
}

// handleUpdate compares the new TimerState with the previous state and sends an update if the state has changed.
func (w *TimerWatch) handleUpdate(s TimerState) {
	if w.lastState != s {
		w.lastState = s
		w.C <- s
	}
}

// Stop stops the TimerWatch.
func (w *TimerWatch) Stop() {
	if w.ticker != nil {
		w.ticker.Stop()
		w.done <- struct{}{}
	}
}
