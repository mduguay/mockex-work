package main

import "fmt"

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	fmt.Println("Hub: Creating new hub")
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	fmt.Println("Hub: Running...")
	for {
		select {
		case client := <-h.register:
			fmt.Println("Hub: Registering connection")
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				fmt.Println("Hub: Unregistering connection")
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
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
