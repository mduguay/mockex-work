package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/mockex", mockexStreamer)
	myRouter.HandleFunc("/holdings/{uid}", holdingHandler)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	fmt.Println("Hello")

	var stg Storage
	stg.connect()
	defer stg.disconnect()

	result := make(chan string)
	go stg.readTrader(1, result)
	fmt.Println(<-result)

	cs := new(CompanyScanner)
	companies := stg.readMultiple(cs)
	for _, c := range companies {
		fmt.Println(*c.(*Company))
	}
	fmt.Printf("%+v\n", companies)
	handleRequests()
}

func holdingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["uid"]

	fmt.Println("Getting holdings for user ", key)
}

func mockexStreamer(w http.ResponseWriter, r *http.Request) {
	hub := newHub()
	// go hub.run()
	// var mkt Market
	// mkt.init()
	// go mkt.openingBell(hub.broadcast)
	// setupEndpoint(hub)
	serveWs(hub, w, r)
}

func setupEndpoint(hub *Hub) {
	http.HandleFunc("/mockex", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
}

func check(err error) {
	if err != nil {
		log.Println(err)
	}
}
