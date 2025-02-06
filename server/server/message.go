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

type VersionMessage struct {
	Version string `msg:"version"`
}

type TimerStateMessage struct {
	Key              string `msg:"K"`
	Active           bool   `msg:"a"`
	Black            bool   `msg:"b"`
	Running          bool   `msg:"r"`
	Countdown        bool   `msg:"c"`
	RemainingSeconds int64  `msg:"s"`
	Clients          int    `msg:"C"`
}

type Command struct {
	Key     string
	Command string `msg:"cmd"`
	Seconds int64  `msg:"sec"`
}

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
