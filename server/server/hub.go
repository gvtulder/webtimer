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
	"log"
	"time"
	"webtimer/server/model"

	"github.com/gorilla/websocket"
)

// check for stale sessions every n seconds
const CLEANUP_INTERVAL = 300

// remove stale sessions after at least n seconds since the last client left
const CLEANUP_DELAY = 600

type Session struct {
	key            string
	timer          *model.Timer
	watch          *model.TimerWatch
	done           chan struct{}
	clients        map[*Client]struct{}
	lastDisconnect int64
}

func newSession(h *Hub, key string) *Session {
	timer := model.NewTimer()
	watch := model.NewTimerWatch(timer)
	session := &Session{key: key, timer: timer, watch: watch, done: make(chan struct{}), clients: make(map[*Client]struct{})}
	h.sessions[key] = session

	watch.Start()
	go func() {
		for {
			select {
			case s := <-watch.C:
				h.broadcast <- TimerStateMessage{
					Key:              key,
					Active:           s.Active,
					Black:            s.Black,
					Countdown:        s.Countdown,
					Running:          s.Running,
					RemainingSeconds: s.Remaining,
					Clients:          len(session.clients),
				}
			case <-session.done:
				return
			}
		}
	}()

	return session
}

func (s *Session) deleteClient(client *Client) {
	delete(s.clients, client)
	if len(s.clients) == 0 {
		s.lastDisconnect = time.Now().Unix()
	}
}

func (s *Session) stop() {
	s.watch.Stop()
	s.done <- struct{}{}
}

type Hub struct {
	logger         *log.Logger
	broadcast      chan TimerStateMessage
	commands       chan Command
	sessions       map[string]*Session
	clients        map[*Client]struct{}
	register       chan *Client
	unregister     chan *Client
	cleanup        *time.Ticker
	welcomeMessage *websocket.PreparedMessage
}

func newHub(logger *log.Logger, version string) *Hub {
	welcomeMessage := VersionMessage{Version: version}
	p, err := welcomeMessage.prepareWebsocketMessage()
	if err != nil {
		logger.Fatalf("Error preparing welcome message: %v", err)
	}

	return &Hub{
		logger:         logger,
		broadcast:      make(chan TimerStateMessage),
		commands:       make(chan Command),
		sessions:       make(map[string]*Session),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		clients:        make(map[*Client]struct{}),
		welcomeMessage: p,
	}
}

func (h *Hub) getSession(key string) *Session {
	session, ok := h.sessions[key]
	if !ok {
		session = newSession(h, key)
		h.sessions[key] = session
	}
	return session
}

func (h *Hub) sendUpdate(session *Session) {
	s := session.timer.State()
	msg := TimerStateMessage{
		Key:              session.key,
		Active:           s.Active,
		Black:            s.Black,
		Countdown:        s.Countdown,
		Running:          s.Running,
		RemainingSeconds: s.Remaining,
		Clients:          len(session.clients),
	}
	h.sendToClients(msg)
}

func (h *Hub) sendToClients(msg TimerStateMessage) {
	p, err := msg.prepareWebsocketMessage()
	if err != nil {
		h.logger.Printf("Error preparing message: %v", err)
		return
	}
	session := h.getSession(msg.Key)
	for client := range session.clients {
		select {
		case client.send <- p:
		default:
			h.logger.Printf("Closing connection [%s]: %s\n", client.key, client.conn.RemoteAddr())
			h.getSession(client.key).deleteClient(client)
			delete(h.clients, client)
			close(client.send)
			h.sendUpdate(session)
		}
	}
}

func (h *Hub) run() {
	h.cleanup = time.NewTicker(CLEANUP_INTERVAL * time.Second)
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
			session := h.getSession(client.key)
			session.clients[client] = struct{}{}
			h.logger.Printf("New connection [%s]: %s (client %d)\n", client.key, client.conn.RemoteAddr(), len(session.clients))
			client.send <- h.welcomeMessage
			h.sendUpdate(session)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.logger.Printf("Closing connection [%s]: %s\n", client.key, client.conn.RemoteAddr())
				session := h.getSession(client.key)
				session.deleteClient(client)
				delete(h.clients, client)
				close(client.send)
				h.sendUpdate(session)
			}

		case msg := <-h.broadcast:
			h.sendToClients(msg)

		case cmd := <-h.commands:
			timer := h.getSession(cmd.Key).timer
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

		case <-h.cleanup.C:
			h.logger.Println("Cleaning time!")
			for key, session := range h.sessions {
				h.logger.Printf("Session %s: %d clients", key, len(session.clients))
				if len(session.clients) == 0 && session.lastDisconnect < time.Now().Unix()-CLEANUP_DELAY {
					h.logger.Println("Closing session " + key)
					session.stop()
					delete(h.sessions, key)
				}
			}
		}
	}
}
