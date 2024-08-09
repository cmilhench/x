package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a single connected client.
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// SocketServer manages client connections and broadcasting.
type SocketServer struct {
	clients   map[*Client]bool
	broadcast chan []byte
	join      chan *Client
	leave     chan *Client
	mu        sync.Mutex
}

// NewSocketServer creates a new SocketServer.
func NewSocketServer() *SocketServer {
	return &SocketServer{
		clients:   make(map[*Client]bool),
		broadcast: make(chan []byte),
		join:      make(chan *Client),
		leave:     make(chan *Client),
	}
}

// Start initializes the socket server to accept new clients and broadcast messages.
func (s *SocketServer) Start() {
	go func() {
		for {
			select {
			case client := <-s.join:
				s.mu.Lock()
				s.clients[client] = true
				s.mu.Unlock()
				log.Printf("Client joined: %v", client.conn.RemoteAddr())

			case client := <-s.leave:
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

// Broadcast sends a message to all connected clients.
func (s *SocketServer) Broadcast(message []byte) {
	s.broadcast <- message
}

// HandleConnections upgrades HTTP connections to WebSocket and handles the client lifecycle.
func (s *SocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte),
	}

	s.join <- client

	go s.handleMessages(client)
}

// handleMessages handles incoming messages from a client.
func (s *SocketServer) handleMessages(client *Client) {
	defer func() {
		s.leave <- client
		client.conn.Close()
	}()

	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		s.Broadcast(msg)
	}
}

// WriteMessages sends messages from the server to the client's WebSocket.
func (client *Client) WriteMessages() {
	for msg := range client.send {
		err := client.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}
}

func main() {
	socketServer := NewSocketServer()
	socketServer.Start()

	http.HandleFunc("/ws", socketServer.HandleConnections)

	log.Println("Socket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
