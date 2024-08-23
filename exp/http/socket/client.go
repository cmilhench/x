package socket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"

	"github.com/cmilhench/x/exp/uuid"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
	id   string
	Name string
}

type MessageHandler func(*Client, []byte)

func NewClient(conn *websocket.Conn) *Client {
	id, _ := uuid.New4()
	return &Client{
		id:   id,
		conn: conn,
		send: make(chan []byte),
	}
}

func (client *Client) ReadMessages(fn MessageHandler) {
	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		fn(client, msg)
	}
}

func (client *Client) WriteMessages() {
	for msg := range client.send {
		err := client.conn.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}
}

func (client *Client) Send(data []byte) {
	client.send <- data
}

func (client *Client) Close() {
	close(client.send)
	deadline := time.Now().Add(5 * time.Second)
	data := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	_ = client.conn.WriteControl(websocket.CloseMessage, data, deadline)
}
