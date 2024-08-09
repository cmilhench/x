package socket

import (
	"log"
	"net/http"
	"sync"

	"github.com/cmilhench/x/exp/uuid"
	"github.com/gorilla/websocket"
)

type SocketServer struct {
	clients   map[*Client]struct{}
	broadcast chan []byte
	join      chan *Client
	part      chan *Client
	mu        sync.Mutex
}

func NewSocketServer() *SocketServer {
	return &SocketServer{
		clients:   make(map[*Client]struct{}),
		broadcast: make(chan []byte),
		join:      make(chan *Client),
		part:      make(chan *Client),
	}
}

func (s *SocketServer) Start() {
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
					delete(s.clients, client)
					close(client.send)
					log.Printf("Client left: %v", client.conn.RemoteAddr())
				}
				s.mu.Unlock()

			case message := <-s.broadcast:
				s.mu.Lock()
				for client := range s.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(s.clients, client)
					}
				}
				s.mu.Unlock()
			}
		}
	}()
}

func (s *SocketServer) Broadcast(message []byte) {
	s.broadcast <- message
}

func (s *SocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	id, _ := uuid.New()
	client := &Client{
		conn:   conn,
		send:   make(chan []byte),
		server: s,
		id:     id,
	}

	s.join <- client

	go client.ReadMessages()
	go client.WriteMessages()
}
