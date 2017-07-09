package gsub

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type LogMessage struct {
	Level   int
	Message string
	uid     string
}

type Logger struct {
	Broker
	msgQ chan LogMessage
}

type ClientInfo struct {
	found bool
	uid   string // Assign a unique ID for each client
}

// A single Broker will be created in this program. It is responsible
// for keeping a list of which clients (browsers) are currently attached
// and broadcasting events (messages) to those clients.
//
type Broker struct {

	// Create a map of clients, the keys of the map are the channels
	// over which we can push messages to attached clients.  (The values
	// are just booleans and are meaningless.)
	//
	clients map[chan LogMessage]ClientInfo

	// Channel into which new clients can be pushed
	//
	newClients chan chan LogMessage

	// Channel into which disconnected clients should be pushed
	//
	defunctClients chan chan LogMessage

	// // Channel into which messages are pushed to be broadcast out
	// // to attahed clients.
	// //
	// messages chan string
}

func (l Logger) Queue() chan LogMessage {
	return l.msgQ
}

func (l *Logger) Start(prefix string) {

	// Make a new Broker instance
	l.Broker = Broker{
		make(map[chan LogMessage]ClientInfo),
		make(chan (chan LogMessage)),
		make(chan (chan LogMessage)),
		// make(chan string),
	}

	l.msgQ = make(chan LogMessage)

	// Start processing events
	l.StartService()

	// Make b the HTTP handler for "/events/".  It can do
	// this because it has a ServeHTTP method.  That method
	// is called in a separate goroutine for each
	// request to "/events/".
	// http.Handle("/events/", b)
	// ssk router handles it
	router.PathPrefix(prefix).HandlerFunc(l.ServeHTTP)

	// Generate a constant stream of events that get pushed
	// into the Broker's messages channel and are then broadcast
	// out to any clients that are attached.

	// Comes from external source.. into l.msgQ

	// go func() {
	// 	for i := 0; ; i++ {

	// 		// Create a little message to send to clients,
	// 		// including the current time.
	// 		b.messages <- fmt.Sprintf("%d - the time is %v", i, time.Now())

	// 		// Print a nice log message and sleep for 5s.
	// 		log.Printf("Sent message %d ", i)
	// 		time.Sleep(5 * 1e9)

	// 	}
	// }()

}

func (l *Logger) Pusher(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I will push logs here")
}

// This Broker method starts a new goroutine.  It handles
// the addition & removal of clients, as well as the broadcasting
// of messages out to clients that are currently attached.
//
func (b *Logger) StartService() {

	// Start a goroutine
	//
	go func() {

		// Loop endlessly
		//
		for {

			// Block until we receive from one of the
			// three following channels.
			select {

			case s := <-b.newClients:

				// There is a new client attached and we
				// want to start sending them messages.

				b.clients[s] = ClientInfo{true, strconv.FormatInt(int64(rand.Intn(256)), 10)}
				log.Println("Added new client")

			case s := <-b.defunctClients:

				// A client has dettached and we want to
				// stop sending them messages.
				delete(b.clients, s)
				close(s)

				log.Println("Removed client : ", s)

			case msg := <-b.msgQ:

				// There is a new message to send.  For each
				// attached client, push the new message
				// into the client's message channel.
				for s, info := range b.clients {
					msg.uid = info.uid
					s <- msg
				}
				log.Printf("Broadcast message to %d clients", len(b.clients))
			}
		}
	}()
}

// This Broker method handles and HTTP request at the "/events/" URL.
//
func (b *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Make sure that the writer supports flushing.
	//
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Create a new channel, over which the broker can
	// send this client messages.
	messageChan := make(chan LogMessage)

	// Add this client to the map of those that should
	// receive updates
	b.newClients <- messageChan

	// Listen to the closing of the http connection via the CloseNotifier
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		// Remove this client from the map of attached clients
		// when `EventHandler` exits.
		b.defunctClients <- messageChan
		log.Println("HTTP connection just closed.")
	}()

	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Don't close the connection, instead loop 10 times,
	// sending messages and flushing the response each time
	// there is a new message to send along.
	//
	// NOTE: we could loop endlessly; however, then you
	// could not easily detect clients that dettach and the
	// server would continue to send them messages long after
	// they're gone due to the "keep-alive" header.  One of
	// the nifty aspects of SSE is that clients automatically
	// reconnect when they lose their connection.
	//
	// A better way to do this is to use the CloseNotifier
	// interface that will appear in future releases of
	// Go (this is written as of 1.0.3):
	// https://code.google.com/p/go/source/detail?name=3292433291b2
	//
	//  fmt.Fprintf(w, "")
	for {

		// Read from our messageChan.
		msg, open := <-messageChan

		// msg, open := <-messageChan

		if !open {
			// If our messageChan was closed, this means that the client has
			// disconnected.
			break
		}

		// Write to the ResponseWriter, `w`.
		// the starting text "data:" is needed for JavaSSE Event Handling processing else anything can be used  https://www.html5rocks.com/en/tutorials/eventsource/basics/
		fmt.Fprintf(w, "data: Hey %v: [%v ]\n\n", msg.uid, msg.Message)

		// Flush the response.  This is only possible if
		// the repsonse supports streaming.
		f.Flush()
	}

	// Done.
	log.Println("Finished HTTP request at ", r.URL.Path)
}
