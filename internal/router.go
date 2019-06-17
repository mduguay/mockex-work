package internal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	Stg Storage
}

func newRouter(storage Storage) {

}

func (rtr *Router) HandleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/mockex", mockexStreamer)
	myRouter.HandleFunc("/holdings/{uid}", rtr.holdingHandler)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func (rtr *Router) holdingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["uid"]

	result := make(chan string)
	go rtr.Stg.readTrader(1, result)
	fmt.Println(<-result)

	cs := new(CompanyScanner)
	companies := rtr.Stg.readMultiple(cs)
	for _, c := range companies {
		fmt.Println(*c.(*Company))
	}
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
