package server

import (
	"github.com/gorilla/websocket"
)

//go:generate msgp

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
