package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"toussaint/backend/structs"
)

// TODO: how to stop broker completely and gracefully?

type NotifierMessage struct {
	Destination ClientType
	Message     []byte
}

type NotifierClient struct {
	Destination ClientType
	Channel     chan []byte
}

type Notifier struct {
	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan NotifierMessage
	// New client connections
	newClients chan NotifierClient
	// Closed client connections
	closingClients chan NotifierClient
	// Client connections registry
	clients map[ClientType]chan []byte
}

func NewNotifier() *Notifier {
	// Instantiate a broker
	notif := &Notifier{
		Notifier:       make(chan NotifierMessage),
		newClients:     make(chan NotifierClient),
		closingClients: make(chan NotifierClient),
		clients:        make(map[ClientType]chan []byte),
	}
	// Set it running - listening and broadcasting events
	go notif.listen()
	return notif
}

func (notifier *Notifier) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	log.Printf("notifier: got new connection")

	//TODO: parse query parameter to get client type

	// Make sure that the writer supports flushing.
	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	// Each connection registers its own message channel with the Broker's connections registry

	client := NotifierClient{
		Destination: Telegram,
		Channel:     make(chan []byte),
	}
	// Signal the broker that we have a new connection
	notifier.newClients <- client
	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		notifier.closingClients <- client
	}()
	// Listen to connection close and un-register messageChan
	notify := rw.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		notifier.closingClients <- client
	}()
	rw.WriteHeader(http.StatusOK)
	for {
		// Write to the ResponseWriter
		// Server Sent Events compatible
		data := <-client.Channel
		log.Printf("ServeHTTP: received data: %s", string(data))
		fmt.Fprint(rw, string(data))
		// Flush the data immediatly instead of buffering it for later.
		flusher.Flush()
	}
	log.Printf("finished notifier.ServeHTTP")
}

func (notifier *Notifier) listen() {
	log.Printf("started notifier.listen()")
	for {
		select {
		case s := <-notifier.newClients:
			// A new client has connected.
			// Register their message channel
			notifier.clients[s.Destination] = s.Channel
			log.Printf("Client added. %d registered clients", len(notifier.clients))
		case s := <-notifier.closingClients:
			// A client has dettached and we want to
			// stop sending them messages.
			delete(notifier.clients, s.Destination)
			log.Printf("Removed client. %d registered clients", len(notifier.clients))
		case event := <-notifier.Notifier:
			notifier.clients[event.Destination] <- event.Message
			log.Printf("Sent message to %+v: %s", event.Destination, string(event.Message))
		}
	}
}

func (notifier *Notifier) NotifyUser(clientType ClientType, notif structs.UserNotification) error {

	data, err := json.Marshal(notif)
	if err != nil {
		return err
	}

	notifier.Notifier <- NotifierMessage{
		Destination: clientType,
		Message:     data,
	}
	return nil
}
