package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cmilhench/x/exp/http/socket"
	"github.com/cmilhench/x/exp/http/static"
	"github.com/cmilhench/x/exp/irc"
)

//go:embed static/*
var fs embed.FS

func main() {
	server := socket.NewSocketServer()
	server.Handle(socketHandler(server))
	server.Start()

	http.Handle("/", http.FileServer(static.Neutered{Prefix: "static", FileSystem: http.FS(fs)}))
	http.HandleFunc("/ws", server.HandleConnections)

	log.Println("Socket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}

func socketHandler(server *socket.Server) socket.MessageHandler {
	return func(client *socket.Client, messageBytes []byte) {
		message := irc.ParseMessage(string(messageBytes))
		log.Printf("message ->: %#v", message)
		switch message.Command {
		case "INFO": // returns information about the <target> server
			client.Send([]byte(fmt.Sprintf("INFO %s", "This is an IRC server.")))
		case "MOTD": // returns the message of the day
			client.Send([]byte(fmt.Sprintf("MOTD %s", "Welcome to the IRC server!")))
		case "NICK": // allows a client to change their IRC nickname.
			client.Name = message.Params
		case "PING": // tests the presence of a connection
			client.Send([]byte(fmt.Sprintf("PONG %s", message.Params)))
		case "NOTICE", "PRIVMSG": // Sends <message> to <target>, which is usually a user or channel.
			if message.Params[0] == '#' {
				server.Broadcast([]byte(fmt.Sprintf(":%s PRIVMSG %s :%s", client.Name, message.Params, message.Trailing)))
			} else {
				server.Send(message.Params, []byte(fmt.Sprintf(":%s PRIVMSG %s :%s", client.Name, message.Params, message.Trailing)))
			}
		case "QUIT": // disconnects the user from the server.
			server.Part(client)
		case "TIME": // returns the current time on the server
			client.Send([]byte(time.Now().Format(time.RFC1123Z)))
		case "TOPIC": // sets the topic of <channel> to <topic>
		default:
			log.Printf("Unknown message type: %#v", message)
		}
	}
}
