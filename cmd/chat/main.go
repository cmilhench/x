package main

import (
	"log"
	"net/http"

	"github.com/cmilhench/x/cmd/chat/socket"
	"github.com/cmilhench/x/cmd/chat/static"
)

func main() {
	server := socket.NewSocketServer()
	server.Start()

	http.Handle("/", static.Handler())
	http.HandleFunc("/ws", server.HandleConnections)

	log.Println("Socket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
