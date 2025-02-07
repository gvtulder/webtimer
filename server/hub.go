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

	"github.com/gorilla/websocket"
)

// check for stale sessions every n seconds
const CLEANUP_INTERVAL = 300

// remove stale sessions after at least n seconds since the last client left
const CLEANUP_DELAY = 600

// A Session represents a timer session with a Timer followed by one or more Clients.
//
// The session runs a goroutine that broadcasts timer events to the clients.
type Session struct {
	key            string               // the key identifying the timer
	timer          *Timer               // the Timer
	watch          *TimerWatch          // the TimerWatch for the Timer
	done           chan struct{}        // channel to stop the Session's goroutine
	clients        map[*Client]struct{} // the clients connected to this timer
	lastDisconnect int64                // the time when the last remaining client disconnected
}

func newSession(h *Hub, key string) *Session {
	timer := NewTimer()
	watch := NewTimerWatch(timer)
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

// deleteClient removes a client from this session.
func (s *Session) deleteClient(client *Client) {
	delete(s.clients, client)
	if len(s.clients) == 0 {
		s.lastDisconnect = time.Now().Unix()
	}
}

// stop stops the TimerWatch and ends the session's goroutine.
func (s *Session) stop() {
	s.watch.Stop()
	s.done <- struct{}{}
}

// The Hub manages the list of timer Sessions and the Clients connected to the WebSocket.
type Hub struct {
	logger         *log.Logger                // the server's logger
	broadcast      chan TimerStateMessage     // channel for messages to be sent to clients
	commands       chan CommandMessage        // channel for commands sent by the clients
	sessions       map[string]*Session        // set of running sessions mapped by key
	clients        map[*Client]struct{}       // set of connected clients
	register       chan *Client               // channel where new clients are registered
	unregister     chan *Client               // channel where disconnecting clients are unregistered
	cleanup        *time.Ticker               // running a regular cleanup process to remove stale timers
	welcomeMessage *websocket.PreparedMessage // the welcome message sent to all new clients
}

// newHub creates a new Hub to manage clients and timer sessions.
func newHub(logger *log.Logger, version string) *Hub {
	welcomeMessage := VersionMessage{Version: version}
	p, err := welcomeMessage.prepareWebsocketMessage()
	if err != nil {
		logger.Fatalf("Error preparing welcome message: %v", err)
	}

	return &Hub{
		logger:         logger,
		broadcast:      make(chan TimerStateMessage),
		commands:       make(chan CommandMessage),
		sessions:       make(map[string]*Session),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		clients:        make(map[*Client]struct{}),
		welcomeMessage: p,
	}
}

// getSession returns or creates the session for the given key.
func (h *Hub) getSession(key string) *Session {
	session, ok := h.sessions[key]
	if !ok {
		session = newSession(h, key)
		h.sessions[key] = session
	}
	return session
}

// sendUpdate composes a TimerStateMessage and queues it for the clients.
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

// sendToClients sends the TimerStateMessage to the clients observing the timer.
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

// run starts the Hub. The Hub processes client connections/disconnections and timer events.
func (h *Hub) run() {
	h.cleanup = time.NewTicker(CLEANUP_INTERVAL * time.Second)
	for {
		select {
		// handle connecting clients
		case client := <-h.register:
			h.clients[client] = struct{}{}
			session := h.getSession(client.key)
			session.clients[client] = struct{}{}
			h.logger.Printf("New connection [%s]: %s (client %d)\n", client.key, client.conn.RemoteAddr(), len(session.clients))
			client.send <- h.welcomeMessage
			h.sendUpdate(session)

		// handle disconnecting clients
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.logger.Printf("Closing connection [%s]: %s\n", client.key, client.conn.RemoteAddr())
				session := h.getSession(client.key)
				session.deleteClient(client)
				delete(h.clients, client)
				close(client.send)
				h.sendUpdate(session)
			}

		// handle timer state updates
		case msg := <-h.broadcast:
			h.sendToClients(msg)

		// handle commands from clients
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

		// periodically remove sessions without active clients
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
