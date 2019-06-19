package internal

import (
	"encoding/json"
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
	myRouter.HandleFunc("/quotes", rtr.quoteHandler)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func (rtr *Router) quoteHandler(w http.ResponseWriter, r *http.Request) {
	cs := new(QuoteScanner)
	companies := rtr.Storage.readMultiple(cs)
	json.NewEncoder(w).Encode(companies)
}

func (rtr *Router) holdingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["uid"]
	hs := new(HoldingScanner)
	hs.uid = key
	holdings := rtr.Storage.readMultiple(hs)
	json.NewEncoder(w).Encode(holdings)
}

func (rtr *Router) traderHandler(w http.ResponseWriter, r *http.Request) {
	result := make(chan string)
	go rtr.Storage.readTrader(1, result)
	fmt.Println(<-result)
}

func (rtr *Router) mockexStreamer(w http.ResponseWriter, r *http.Request) {
	serveWs(rtr.Hub, w, r)
}
