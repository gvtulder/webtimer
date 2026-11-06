// Copyright (C) 2026 Gijs van Tulder
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
	"encoding/json"
	"log"
	"os"
	"sync"
)

// A PersistedTimer is used to persist the timer state to JSON.
type PersistedTimer struct {
	Black            bool  `json:"black"`            // true when in black-screen mode
	Countdown        bool  `json:"countdown"`        // true when in countdown mode
	StartTime        int64 `json:"startTime"`        // if the timer is running, the time when when the current run started, otherwise 0
	RemainingAtStart int64 `json:"remainingAtStart"` // the milliseconds-on-the-clock at startTime
}

// A PersistedTimer is used to persist the timer state to JSON.
type PersistedTimers struct {
	Timers map[string]PersistedTimer `json:"timers"` // a map of active timers
}

// A PersistenceLoop is responsible for the asynchronous saving of the timer states to a JSON file.
type PersistenceLoop struct {
	updates chan PersistedTimers // channel to trigger state saving
	done    chan struct{}
	logger  *log.Logger
	path    string

	mu     sync.Mutex
	closed bool
}

// Start a new PersistenceLoop.
func NewPersistenceLoop(logger *log.Logger, path string) *PersistenceLoop {
	p := &PersistenceLoop{
		updates: make(chan PersistedTimers, 1),
		done:    make(chan struct{}),
		logger:  logger,
		path:    path,
	}
	go func() {
		defer close(p.done)
		for state := range p.updates {
			p.write(state)
		}
	}()
	return p
}

// write writes the state to the file.
func (p *PersistenceLoop) write(state PersistedTimers) error {
	data, err := json.Marshal(state)
	if err != nil {
		p.logger.Printf("Could not persist state.")
		return err
	}
	tmp := p.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		p.logger.Printf("Could not persist state to %s: %v", p.path, err)
		return err
	}
	if err := os.Rename(tmp, p.path); err != nil {
		p.logger.Printf("Could not persist state to %s: %v", p.path, err)
		return err
	}
	return nil
}

// Save queues the state for saving to the state file.
//
// The state is written asynchronously. If multiple updates are received
// quickly after another, only the most recent state will be written.
func (p *PersistenceLoop) Save(state PersistedTimers) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return
	}

	select {
	case p.updates <- state:
	default:
		select {
		case <-p.updates: // discard the older snapshot
		default:
		}
		p.updates <- state
	}
}

// Close ends the PersistenceLoop.
func (p *PersistenceLoop) Close() {
	p.mu.Lock()
	if !p.closed {
		p.closed = true
		close(p.updates)
	}
	p.mu.Unlock()
	<-p.done
}

// LoadTimerState loads the timer states from a JSON file.
func LoadTimerState(path string) (*PersistedTimers, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	state := &PersistedTimers{}
	if err := json.Unmarshal(data, state); err != nil {
		return nil, err
	}
	return state, nil
}
