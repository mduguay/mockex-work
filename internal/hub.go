package internal

import "log"

type Hub struct {
	clients    map[*Client]bool
	Broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	log.Println("Hub: Creating new hub")
	return &Hub{
		Broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	log.Println("Hub: Running...")
	for {
		select {
		case client := <-h.register:
			log.Println("Hub: Registering connection")
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Println("Hub: Unregistering connection")
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.Broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
					log.Println("Hub: Sent message")
				default:
					log.Println("Hub: Closing connection")
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
