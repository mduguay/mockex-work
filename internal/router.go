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
	myRouter.HandleFunc("/companies", rtr.companyHandler)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func (rtr *Router) companyHandler(w http.ResponseWriter, r *http.Request) {
	cs := new(CompanyScanner)
	companies := rtr.Stg.readMultiple(cs)
	for _, c := range companies {
		fmt.Println(*c.(*Company))
	}
}

func (rtr *Router) holdingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["uid"]

	fmt.Println("Getting holdings for user ", key)

	hs := new(HoldingScanner)
	hs.uid = key
	holdings := rtr.Stg.readMultiple(hs)
	for _, c := range holdings {
		fmt.Println(*c.(*Holding))
	}
}

func (rtr *Router) traderHandler(w http.ResponseWriter, r *http.Request) {
	result := make(chan string)
	go rtr.Stg.readTrader(1, result)
	fmt.Println(<-result)
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
