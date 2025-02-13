package stream

import (
	"fmt"
	"net/http"

	"github.com/cmilhench/x/exp/uuid"
)

func NewEventSourceClient(w http.ResponseWriter, r *http.Request) (*EventSourceClient, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("streaming unsupported")
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	id, _ := uuid.New4()
	return &EventSourceClient{
		id:      id,
		send:    make(chan []byte),
		w:       w,
		r:       r,
		flusher: flusher,
	}, nil
}

type EventSourceClient struct {
	send    chan []byte
	id      string
	w       http.ResponseWriter
	r       *http.Request
	flusher http.Flusher
}

func (client *EventSourceClient) WriteMessages() {
	for {
		select {
		case message := <-client.send:
			// Here we send only the "data" field, but there are few others
			fmt.Fprintf(client.w, "event: message\n")
			fmt.Fprintf(client.w, "data: %s\n\n", message)
			client.flusher.Flush()
		case <-client.r.Context().Done():
			return
		}
	}
}

func (client *EventSourceClient) Send(data []byte) bool {
	select {
	case client.send <- data:
		return true
	default:
		close(client.send)
		return false
	}
}

func (client *EventSourceClient) Identifier() string {
	return client.id
}

func (client *EventSourceClient) Close() {
	close(client.send)
}
