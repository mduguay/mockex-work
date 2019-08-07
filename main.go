package main

import (
	"log"

	"github.com/mduguay/mockex-work/internal"
)

func main() {
	log.Println("Init: Storage")
	sto := new(internal.Storage)
	sto.Connect()
	defer sto.Disconnect()

	log.Println("Init: Hub")
	hub := internal.NewHub()
	go hub.Run()

	log.Println("Init: Market")
	mkt := internal.Market{
		Storage: sto,
	}
	go mkt.OpeningBell(hub.Broadcast)

	log.Println("Init: Router")
	router := internal.Router{
		Storage: sto,
		Hub:     hub,
		Market:  &mkt,
	}
	router.HandleRequests()
}
