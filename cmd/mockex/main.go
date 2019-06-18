package main

import (
	"git.csnzoo.com/mduguay/mockex/internal"
)

func main() {
	sto := new(internal.Storage)
	sto.Connect()
	defer sto.Disconnect()

	hub := internal.NewHub()
	go hub.Run()

	mkt := internal.NewMarket()
	go mkt.OpeningBell(hub.Broadcast)

	router := internal.Router{
		Storage: sto,
		Hub:     hub,
	}
	router.HandleRequests()
}
