package server

import (
	"webtimer/server/model"
)

type TimerStateMessage struct {
	Active           bool  `json:"active"`
	Black            bool  `json:"black"`
	Running          bool  `json:"running"`
	Countdown        bool  `json:"countdown"`
	RemainingSeconds int64 `json:"remaining"`
	Clients          int   `json:"clients"`
}

type Command struct {
	Command string `json:"cmd"`
	Seconds int64  `json:"seconds"`
}

type Hub struct {
	broadcast  chan TimerStateMessage
	commands   chan Command
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan TimerStateMessage),
		commands:   make(chan Command),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run(watch *model.TimerWatch, timer *model.Timer) {
	go func() {
		for s := range watch.C {
			h.broadcast <- TimerStateMessage{
				Active:           s.Active,
				Black:            s.Black,
				Countdown:        s.Countdown,
				Running:          s.Running,
				RemainingSeconds: s.Remaining,
				Clients:          len(h.clients),
			}
		}
	}()

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			s := timer.State()
			client.send <- TimerStateMessage{
				Active:           s.Active,
				Black:            s.Black,
				Countdown:        s.Countdown,
				Running:          s.Running,
				RemainingSeconds: s.Remaining,
				Clients:          len(h.clients),
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}

		case cmd := <-h.commands:
			switch cmd.Command {
			case "start":
				timer.Start()
			case "reset":
				timer.Reset()
			case "pause":
				timer.Pause()
			case "set":
				timer.SetRemaining(cmd.Seconds)
			case "add":
				timer.AddRemaining(cmd.Seconds)
			case "blackon":
				timer.BlackOn()
			case "blackoff":
				timer.BlackOff()
			case "toggleblack":
				timer.ToggleBlack()
			}
		}
	}
}
