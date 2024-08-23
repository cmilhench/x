package socket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	clients   map[*Client]struct{}
	broadcast chan []byte
	messages  chan struct {
		Target string
		Data   []byte
	}
	join    chan *Client
	part    chan *Client
	handler MessageHandler
	mu      sync.Mutex
}

func NewSocketServer() *Server {
	return &Server{
		clients:   make(map[*Client]struct{}),
		broadcast: make(chan []byte),
		messages: make(chan struct {
			Target string
			Data   []byte
		}),
		join: make(chan *Client),
		part: make(chan *Client),
	}
}

func (s *Server) Start() {
	go func() {
		for {
			select {
			case client := <-s.join:
				s.mu.Lock()
				s.clients[client] = struct{}{}
				s.mu.Unlock()
				log.Printf("Client joined: %v", client.conn.RemoteAddr())
			case client := <-s.part:
				s.mu.Lock()
				if _, ok := s.clients[client]; ok {
					client.Close()
					delete(s.clients, client)
					log.Printf("Client left: %v", client.conn.RemoteAddr())
				}
				s.mu.Unlock()
			case data := <-s.broadcast:
				s.mu.Lock()
				for client := range s.clients {
					select {
					case client.send <- data:
					default:
						close(client.send)
						delete(s.clients, client)
					}
				}
				s.mu.Unlock()
			case message := <-s.messages:
				s.mu.Lock()
				for k := range s.clients {
					if k.id == message.Target || k.Name == message.Target {
						select {
						case k.send <- message.Data:
						default:
							close(k.send)
							delete(s.clients, k)
						}
						return
					}
				}
				s.mu.Unlock()
			}
		}
	}()
}

func (s *Server) Broadcast(message []byte) {
	s.broadcast <- message
}

func (s *Server) Handle(handler MessageHandler) {
	s.handler = handler
}

func (s *Server) Send(target string, message []byte) {
	s.messages <- struct {
		Target string
		Data   []byte
	}{target, message}
}

func (s *Server) Part(client *Client) {
	s.part <- client
}

func (s *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	client := NewClient(conn)

	s.join <- client

	go client.WriteMessages()
	client.ReadMessages(s.handler)
}
