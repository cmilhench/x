package stream

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cmilhench/x/exp/logger"
	"github.com/cmilhench/x/exp/uuid"

	"github.com/gorilla/websocket"
)

type MessageHandler func(Client, []byte)

type SocketClient struct {
	conn *websocket.Conn
	send chan []byte
	id   string
	Name string
}

func NewSocketClient(w http.ResponseWriter, r *http.Request) (*SocketClient, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, fmt.Errorf("Upgrade error: %w", err)
	}
	id, _ := uuid.New4()
	return &SocketClient{
		id:   id,
		conn: conn,
		send: make(chan []byte),
	}, nil
}

func (client *SocketClient) ReadMessages(fn MessageHandler) {
	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			logger.Debug(fmt.Sprintf("Read error: %v", err))
			break
		}
		fn(client, msg)
	}
}

func (client *SocketClient) WriteMessages() {
	for msg := range client.send {
		err := client.conn.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			logger.Debug(fmt.Sprintf("Write error: %v", err))
			break
		}
	}
}

func (client *SocketClient) Send(data []byte) bool {
	select {
	case client.send <- data:
		return true
	default:
		close(client.send)
		return false
	}
}

func (client *SocketClient) Identifier() string {
	return client.id
}

func (client *SocketClient) Close() {
	close(client.send)
	deadline := time.Now().Add(5 * time.Second)
	data := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	_ = client.conn.WriteControl(websocket.CloseMessage, data, deadline)
}
