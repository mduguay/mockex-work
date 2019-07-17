package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Router struct {
	Storage *Storage
	Hub     *Hub
}

func (rtr *Router) HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login", rtr.loginHandler)
	router.HandleFunc("/mockex", rtr.mockexStreamer)
	router.HandleFunc("/trader/{tid}", rtr.traderHandler)
	router.HandleFunc("/holdings/{tid}", rtr.holdingHandler)
	router.HandleFunc("/quotes", rtr.quoteHandler)
	router.HandleFunc("/trade", rtr.tradeHandler).Methods("POST")
	router.HandleFunc("/cash/{tid}", rtr.cashHandler)
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
	decoder := json.NewDecoder(r.Body)
	var t Trade
	err := decoder.Decode(&t)
	check(err)
	fmt.Println(t)
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

func (rtr *Router) mockexStreamer(w http.ResponseWriter, r *http.Request) {
	serveWs(rtr.Hub, w, r)
}
