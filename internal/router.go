package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Router struct {
	Storage *Storage
	Hub     *Hub
}

func (rtr *Router) HandleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/mockex", rtr.mockexStreamer)
	myRouter.HandleFunc("/trader/{tid}", rtr.traderHandler)
	myRouter.HandleFunc("/holdings/{tid}", rtr.holdingHandler)
	myRouter.HandleFunc("/quotes", rtr.quoteHandler)
	myRouter.HandleFunc("/trade", rtr.tradeHandler)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func (rtr *Router) quoteHandler(w http.ResponseWriter, r *http.Request) {
	cs := new(QuoteScanner)
	companies := rtr.Storage.readMultiple(cs)
	json.NewEncoder(w).Encode(companies)
}

func (rtr *Router) holdingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["tid"]
	hs := new(HoldingScanner)
	hs.uid = key
	holdings := rtr.Storage.readMultiple(hs)
	json.NewEncoder(w).Encode(holdings)
}

func (rtr *Router) traderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["tid"]
	tid, err := strconv.Atoi(key)
	if check(err) {
		return
	}
	t := rtr.Storage.readTrader(tid)
	json.NewEncoder(w).Encode(t)
}

func (rtr *Router) tradeHandler(w http.ResponseWriter, r *http.Request) {
	// request is post with tid, symbol, amount, direction, price
	decoder := json.NewDecoder(r.Body)
	var t Trade
	err := decoder.Decode(&t)
	check(err)
	fmt.Println(t)
	// publish trade to db
	rtr.Storage.createTrade(t)
	// update holding
	// return holding
}

func (rtr *Router) mockexStreamer(w http.ResponseWriter, r *http.Request) {
	serveWs(rtr.Hub, w, r)
}
