package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	dumpDayPrices()
}

func setupEndpoint() {
	hub := newHub()
	go hub.run()
	http.HandleFunc("/mockex", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	fmt.Println("Serving on ws://127.0.0.1:8080/mockex")
	go tickWriter(hub)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func tickWriter(h *Hub) {
	for {
		time.Sleep(time.Second * 2)
		h.broadcast <- []byte("tick")
	}
}
