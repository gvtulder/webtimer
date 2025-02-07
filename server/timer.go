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
	"math"
	"time"
)

// A TimerState describes the state of a timer at a given point in time.
// It can be passed around to other methods, e.g., in event handlers.
//
// In countdown mode (Countdown == true), Remaining represents the remaining available time in seconds.
// This will become negative after the timer runs out.
//
// In countup mode (Countdown == false), Remaining start at 0 and immediately becomes negative.
// To show the number of seconds since starting at 0, use -Remaining.
type TimerState struct {
	Active    bool  // true if the timer is running or shows time != 0
	Black     bool  // true if the timer is in black-screen mode
	Countdown bool  // true if the timer is counting down, false if counting up
	Running   bool  // true if the timer is running
	Remaining int64 // the number of remaining seconds
}

// A Timer is a timer with its corresponding state.
//
// The Timer sends a message to channel C when its state is changed (start/stop/add/set/etc.).
//
// The timer keeps track of the remaining time as follows:
//
//  1. When the timer is started, startTime is set to the time.Now().UnixMilli().
//     remainingAtStart is the number of milliseconds that was remaining at startTime.
//     The remaining time can be computed by remainingAtStart - (time.Now().UnixMilli() - startTime).
//
//  2. If the timer is paused, startTime is set to 0.
//     remainingAtStart is set to the number of milliseconds remaining on the clock.
type Timer struct {
	C                chan struct{} // channel for timer updates
	black            bool          // true when in black-screen mode
	countdown        bool          // true when in countdown mode
	startTime        int64         // if the timer is running, the time when when the current run started, otherwise 0
	remainingAtStart int64         // the milliseconds-on-the-clock at startTime
}

func NewTimer() *Timer {
	return &Timer{C: make(chan struct{})}
}

// Active returns true if the timer is running or shows a time other than 0, false if it has been reset or never been run.
func (t *Timer) Active() bool {
	return t.startTime != 0 || t.remainingAtStart != 0
}

// Black returns true if the timer is currently showing a black screen.
func (t *Timer) Black() bool {
	return t.black
}

// Countdown returns true if the timer is counting down; false if it is counting up.
func (t *Timer) Countdown() bool {
	return t.countdown
}

// Running returns true if the timer is currently running.
func (t *Timer) Running() bool {
	return t.startTime != 0
}

// Remaining returns the number of remaining seconds at the current time.
func (t *Timer) Remaining() int64 {
	var r int64
	if t.startTime != 0 {
		r = t.remainingAtStart - (time.Now().UnixMilli() - t.startTime)
	} else {
		r = t.remainingAtStart
	}
	return int64(math.Round(float64(r) / 1000))
}

// State returns the current state of the timer.
func (t *Timer) State() TimerState {
	return TimerState{
		Active:    t.Active(),
		Black:     t.Black(),
		Countdown: t.Countdown(),
		Running:   t.Running(),
		Remaining: t.Remaining(),
	}
}

// SetRemaining sets the timer to s seconds.
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

// AddRemaining adds s seconds to the timer (s can be negative).
func (t *Timer) AddRemaining(s int64) {
	t.remainingAtStart += s * 1000
	if t.Remaining() > 0 {
		t.countdown = true
	}
	t.C <- struct{}{}
}

// Start starts the timer, if it is not already running.
func (t *Timer) Start() {
	if t.startTime == 0 {
		t.startTime = time.Now().UnixMilli()
		t.C <- struct{}{}
	}
}

// Start pauses the timer, if it is not already running.
func (t *Timer) Pause() {
	if t.Running() {
		t.remainingAtStart = t.remainingAtStart - (time.Now().UnixMilli() - t.startTime)
		t.startTime = 0
		t.C <- struct{}{}
	}
}

// Reset resets the timer to zero.
func (t *Timer) Reset() {
	if t.Active() {
		t.startTime = 0
		t.remainingAtStart = 0
		t.countdown = false
		t.C <- struct{}{}
	}
}

// ToggleBlack toggles to/from black-screen mode.
func (t *Timer) ToggleBlack() {
	t.black = !t.black
	t.C <- struct{}{}
}

// BlackOn enables black-screen mode.
func (t *Timer) BlackOn() {
	if !t.black {
		t.black = true
		t.C <- struct{}{}
	}
}

// BlackOff disables black-screen mode.
func (t *Timer) BlackOff() {
	if t.black {
		t.black = false
		t.C <- struct{}{}
	}
}
