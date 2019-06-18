package internal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	Storage *Storage
	Hub     *Hub
}

func (rtr *Router) HandleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/mockex", rtr.mockexStreamer)
	myRouter.HandleFunc("/holdings/{uid}", rtr.holdingHandler)
	myRouter.HandleFunc("/companies", rtr.companyHandler)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func (rtr *Router) companyHandler(w http.ResponseWriter, r *http.Request) {
	cs := new(CompanyScanner)
	companies := rtr.Storage.readMultiple(cs)
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
	holdings := rtr.Storage.readMultiple(hs)
	for _, c := range holdings {
		fmt.Println(*c.(*Holding))
	}
}

func (rtr *Router) traderHandler(w http.ResponseWriter, r *http.Request) {
	result := make(chan string)
	go rtr.Storage.readTrader(1, result)
	fmt.Println(<-result)
}

func (rtr *Router) mockexStreamer(w http.ResponseWriter, r *http.Request) {
	serveWs(rtr.Hub, w, r)
}
