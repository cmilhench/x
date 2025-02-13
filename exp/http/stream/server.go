package stream

import (
	"net/http"
	"sync"
)

type Message []byte

type DirectMessage struct {
	Target string
	Data   Message
}

type Client interface {
	Send(data []byte) bool
	Identifier() string
	Close()
}

type Server struct {
	clients   map[Client]struct{}
	broadcast chan Message
	messages  chan DirectMessage
	join      chan Client
	part      chan Client
	mu        sync.Mutex
}

func NewServer() *Server {
	return &Server{
		clients:   make(map[Client]struct{}),
		broadcast: make(chan Message),
		messages:  make(chan DirectMessage),
		join:      make(chan Client),
		part:      make(chan Client),
	}
}

func (s *Server) Listen() {
	for {
		select {
		case client := <-s.join:
			s.mu.Lock()
			s.clients[client] = struct{}{}
			s.mu.Unlock()
			//log.Debugf("Client joined: %v", client.Identifier())
		case client := <-s.part:
			s.mu.Lock()
			if _, ok := s.clients[client]; ok {
				client.Close()
				delete(s.clients, client)
				//log.Debugf("Client left: %v", client.Identifier())
			}
			s.mu.Unlock()
		case data := <-s.broadcast:
			s.mu.Lock()
			for client := range s.clients {
				if !client.Send(data) {
					delete(s.clients, client)
				}
			}
			s.mu.Unlock()
		case message := <-s.messages:
			s.mu.Lock()
			for k := range s.clients {
				if k.Identifier() == message.Target {
					if !k.Send(message.Data) {
						delete(s.clients, k)
					}
					return
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *Server) Broadcast(message []byte) {
	s.broadcast <- message
}

func (s *Server) Send(target string, message []byte) {
	s.messages <- DirectMessage{target, message}
}

func (s *Server) Part(client Client) {
	s.part <- client
}

func (s *Server) WebsocketHandler(handler MessageHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client, err := NewSocketClient(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.join <- client

		go client.ReadMessages(handler)
		client.WriteMessages()
	}
}

func (s *Server) EventSourceHandler(w http.ResponseWriter, r *http.Request) {
	client, err := NewEventSourceClient(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.join <- client

	client.WriteMessages()
}
