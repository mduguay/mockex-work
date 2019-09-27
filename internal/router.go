package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mduguay/mockex-work/data"
	"github.com/rs/cors"
)

// Router is the object that handles http requests
type Router struct {
	Storage *Storage
	Hub     *Hub
	Market  *Market
}

// HandleRequests captures all incoming http requests and routes them accordingly
func (rtr *Router) HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/mockex", rtr.mockexStreamer)
	router.HandleFunc("/login", rtr.loginHandler)
	router.HandleFunc("/trader/{tid}", rtr.traderHandler)
	router.HandleFunc("/holdings/{tid}", rtr.holdingHandler)
	router.HandleFunc("/quotes", rtr.quoteHandler)
	router.HandleFunc("/trade", rtr.tradeHandler).Methods("POST")
	router.HandleFunc("/cash/{tid}", rtr.cashHandler)
	router.HandleFunc("/history/{cid}", rtr.historyHandler)
	router.HandleFunc("/settings/{cid}", rtr.settingsHandler).Methods("POST")
	router.HandleFunc("/market/{action}", rtr.marketHandler)
	router.HandleFunc("/backfill", rtr.backfillHandler)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("ALLOWED_ORIGIN")},
		AllowCredentials: true,
	})
	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), handler))
}

func (rtr *Router) loginHandler(w http.ResponseWriter, r *http.Request) {
	t := rtr.Storage.readTrader(1)
	json.NewEncoder(w).Encode(t)
}

func (rtr *Router) quoteHandler(w http.ResponseWriter, r *http.Request) {
	cs := new(data.QuoteScanner)
	companies := rtr.Storage.readMultiple(cs)
	json.NewEncoder(w).Encode(companies)
}

func (rtr *Router) holdingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["tid"]
	hs := new(data.HoldingScanner)
	hs.UID = key
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
	decoder := json.NewDecoder(r.Body)
	t := new(Trade)
	err := decoder.Decode(t)
	check(err)
	log.Printf("Trade: %+v", t)
	shares, err := rtr.Storage.createTrade(t)
	check(err)
	json.NewEncoder(w).Encode(shares)
}

func (rtr *Router) cashHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["tid"]
	tid, err := strconv.Atoi(key)
	if check(err) {
		return
	}
	c := rtr.Storage.readCash(tid)
	json.NewEncoder(w).Encode(c)
}

func (rtr *Router) historyHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["cid"]
	k, err := strconv.Atoi(key)
	if check(err) {
		return
	}
	hs := new(data.HistoryScanner)
	hs.Cid = k
	holdings := rtr.Storage.readMultiple(hs)
	json.NewEncoder(w).Encode(holdings)
}

func (rtr *Router) settingsHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["cid"]
	k, err := strconv.Atoi(key)
	if check(err) {
		return
	}
	decoder := json.NewDecoder(r.Body)
	s := new(Settings)
	err = decoder.Decode(s)
	check(err)
	rtr.Storage.updateSettings(k, s)
}

func (rtr *Router) marketHandler(w http.ResponseWriter, r *http.Request) {
	action := mux.Vars(r)["action"]
	fmt.Println("Router: Market: Action:", action)
	switch action {
	case "start":
		rtr.Market.OpeningBell(rtr.Hub.Broadcast)
	case "stop":
		rtr.Market.ClosingBell()
	}
}

func (rtr *Router) backfillHandler(w http.ResponseWriter, r *http.Request) {
	rtr.Market.backfill()
	return
}

func (rtr *Router) mockexStreamer(w http.ResponseWriter, r *http.Request) {
	serveWs(rtr.Hub, w, r)
}
