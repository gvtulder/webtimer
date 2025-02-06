package server

import (
	"log"
	"time"
	"webtimer/server/model"
)

// check for stale sessions every n seconds
const CLEANUP_INTERVAL = 300

// remove stale sessions after at least n seconds since the last client left
const CLEANUP_DELAY = 600

type TimerStateMessage struct {
	Key              string `json:"key"`
	Active           bool   `json:"active"`
	Black            bool   `json:"black"`
	Running          bool   `json:"running"`
	Countdown        bool   `json:"countdown"`
	RemainingSeconds int64  `json:"remaining"`
	Clients          int    `json:"clients"`
}

type Command struct {
	Key     string
	Command string `json:"cmd"`
	Seconds int64  `json:"seconds"`
}

type Session struct {
	timer          *model.Timer
	watch          *model.TimerWatch
	done           chan bool
	clients        map[*Client]bool
	lastDisconnect int64
}

func newSession(h *Hub, key string) *Session {
	timer := model.NewTimer()
	watch := model.NewTimerWatch(timer)
	session := &Session{timer: timer, watch: watch, done: make(chan bool), clients: make(map[*Client]bool)}
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
	s.done <- true
}

type Hub struct {
	broadcast  chan TimerStateMessage
	commands   chan Command
	sessions   map[string]*Session
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	cleanup    *time.Ticker
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan TimerStateMessage),
		commands:   make(chan Command),
		sessions:   make(map[string]*Session),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
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

func (h *Hub) run(logger *log.Logger) {
	h.cleanup = time.NewTicker(CLEANUP_INTERVAL * time.Second)
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			session := h.getSession(client.key)
			session.clients[client] = true
			logger.Printf("New connection [%s]: %s (client %d)\n", client.key, client.conn.RemoteAddr(), len(session.clients))
			s := session.timer.State()
			client.send <- TimerStateMessage{
				Key:              client.key,
				Active:           s.Active,
				Black:            s.Black,
				Countdown:        s.Countdown,
				Running:          s.Running,
				RemainingSeconds: s.Remaining,
				Clients:          len(session.clients),
			}

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				logger.Printf("Closing connection [%s]: %s\n", client.key, client.conn.RemoteAddr())
				h.getSession(client.key).deleteClient(client)
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			session := h.getSession(message.Key)
			for client := range session.clients {
				select {
				case client.send <- message:
				default:
					logger.Printf("Closing connection [%s]: %s\n", client.key, client.conn.RemoteAddr())
					h.getSession(client.key).deleteClient(client)
					delete(h.clients, client)
					close(client.send)
				}
			}

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
			logger.Println("Cleaning time!")
			for key, session := range h.sessions {
				logger.Printf("Session %s: %d clients", key, len(session.clients))
				if len(session.clients) == 0 && session.lastDisconnect < time.Now().Unix()-CLEANUP_DELAY {
					logger.Println("Closing session " + key)
					session.stop()
					delete(h.sessions, key)
				}
			}
		}
	}
}
