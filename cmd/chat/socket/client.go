package socket

import (
	"fmt"
	"log"

	"github.com/cmilhench/x/cmd/chat/irc"
	"github.com/gorilla/websocket"
)

type Client struct {
	server *SocketServer
	conn   *websocket.Conn
	send   chan []byte
	id     string
}

func (client *Client) ReadMessages() {
	defer func() {
		client.server.part <- client
		client.conn.Close()
	}()

	for {
		_, messageBytes, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		var message irc.Message
		message.Parse(string(messageBytes))

		switch message.Command {
		case "PRIVMSG":
			// TODO: Handle PRIVMSG to #channel or to @user
			log.Printf("Received message: %s", message.Trailing)
			client.server.Broadcast(messageBytes)
		case "JOIN":
			// broadcast a message from the server to the #channel
			response := irc.Message{
				Prefix:   "server",
				Command:  "PRIVMSG",
				Params:   message.Params,
				Trailing: fmt.Sprintf("%s joined %s", message.Prefix, message.Params),
			}
			client.server.Broadcast([]byte(response.String()))
		case "something":
			// TODO: send a message only to a given client
			client.send <- []byte("bye")
		default:
			log.Printf("Unknown message type: %#v", message)
		}
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
