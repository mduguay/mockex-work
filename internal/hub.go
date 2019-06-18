package internal

import "fmt"

type Hub struct {
	clients    map[*Client]bool
	Broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	fmt.Println("Hub: Creating new hub")
	return &Hub{
		Broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	fmt.Println("Hub: Running...")
	for {
		select {
		case client := <-h.register:
			fmt.Println("Hub: Registering connection")
			h.clients[client] = true
			client.send <- initialData()
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				fmt.Println("Hub: Unregistering connection")
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.Broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
					fmt.Println("Hub: Sent message")
				default:
					fmt.Println("Hub: Closing connection")
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
