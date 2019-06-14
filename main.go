package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Hello")
	var stg Storage
	stg.connect()
	result := make(chan string)
	go stg.readUser(1, result)
	fmt.Println(<-result)
	stg.disconnect()
	// hub := newHub()
	// go hub.run()
	// var mkt Market
	// mkt.init()
	// go mkt.openingBell(hub.broadcast)
	// setupEndpoint(hub)
}

func setupEndpoint(hub *Hub) {
	http.HandleFunc("/mockex", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	fmt.Println("Serving on ws://127.0.0.1:8080/mockex")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
