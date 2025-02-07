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
	"github.com/gorilla/websocket"
)

//go:generate msgp

// A VersionMessage (in msgpack format) is sent to the client to communicate the server version.
type VersionMessage struct {
	Version string `msg:"version"` // the server version string
}

// A TimerStateMessage (in msgpack format) is sent to the client to communicate the current state of a timer.
type TimerStateMessage struct {
	Key              string `msg:"K"` // the key identifying the timer
	Active           bool   `msg:"a"` // true if the timer is running or shows time != 0
	Black            bool   `msg:"b"` // true if the timer is in black-screen mode
	Countdown        bool   `msg:"c"` // true if the timer is counting down, false if counting up
	Running          bool   `msg:"r"` // true if the timer is running
	RemainingSeconds int64  `msg:"s"` // the number of remaining seconds
	Clients          int    `msg:"C"` // the number of clients connected to this timer
}

// A CommandMessage (in msgpack format) is sent by the client to give commands.
type CommandMessage struct {
	Key     string // the key identifying the timer
	Command string `msg:"cmd"` // the command to be executed
	Seconds int64  `msg:"sec"` // if applicable, the seconds parameter for the command
}

// prepareWebsocketMessage converts a VersionMessage into a websocket.PreparedMessage that can be sent over a websocket.
func (m *VersionMessage) prepareWebsocketMessage() (*websocket.PreparedMessage, error) {
	b, err := m.MarshalMsg(nil)
	if err != nil {
		return nil, err
	}
	p, err := websocket.NewPreparedMessage(websocket.BinaryMessage, b)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// prepareWebsocketMessage converts a TimerStateMessage into a websocket.PreparedMessage that can be sent over a websocket.
func (m *TimerStateMessage) prepareWebsocketMessage() (*websocket.PreparedMessage, error) {
	b, err := m.MarshalMsg(nil)
	if err != nil {
		return nil, err
	}
	p, err := websocket.NewPreparedMessage(websocket.BinaryMessage, b)
	if err != nil {
		return nil, err
	}
	return p, nil
}
