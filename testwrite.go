package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func mockexEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected")
	err = ws.WriteMessage(1, []byte("Hi Client"))
	if err != nil {
		log.Println(err)
	}

	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

		if err := conn.WriteMessage(1, []byte("Hi from Mock Exchange")); err != nil {
			log.Println(err)
			return
		}
	}
}

func setupRoutes() {
	http.HandleFunc("/mockex", mockexEndpoint)
}
