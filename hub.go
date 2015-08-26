// Handles registering / unregistering / sending messages to connected clients

package main

import log "github.com/Sirupsen/logrus"

// Handles sending messages to all connected clients
type hub struct {
	// Stores all connected clients
	connections map[*connection]bool

	// Messages to broadcast to the connected clients
	broadcast chan []byte

	// Register a new client
	register chan *connection

	// Unregister clients, removing them from the pool
	unregister chan *connection
}

// Listens on the channels for messages and performs the
// relivant functionality
func (h *hub) run() {
	log.Info("Starting Hub")
	for {
		select {
		// On register messages, store the connection
		case c := <-h.register:
			log.Debug("Client Registered: ", c.ws.RemoteAddr())
			h.connections[c] = true
		// On unregister messsages, delete the connection from the pool
		case c := <-h.unregister:
			log.Debug("Client Unregistered: ", c.ws.RemoteAddr())
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		// On incoming messages, loop over connected clients and send
		// the message
		case m := <-h.broadcast:
			log.Debug("Broadcast Message: ", string(m[:]))
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					log.Debug("Client no longer active: ", c.ws.RemoteAddr())
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}

// Create a new instance of the hub
var h = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}
