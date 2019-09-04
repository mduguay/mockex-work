package internal

import "log"

// Hub is the collector of data and will publish it to all connected clients
type Hub struct {
	clients    map[*Client]bool
	Broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// NewHub creates and initializes the hub
func NewHub() *Hub {
	log.Println("Hub: Creating new hub")
	return &Hub{
		Broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run starts the hub publishing loop
func (h *Hub) Run() {
	log.Println("Hub: Running...")
	for {
		select {
		case client := <-h.register:
			log.Println("Hub: Registering connection")
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Println("Hub: Closing connection")
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.Broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// In what case does this get hit?
					log.Println("Hub: Closing connection on Broadcast")
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
