package main

import (
	"log"

	"git.csnzoo.com/mduguay/mockex/internal"
)

func main() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatalln("Failed to load environment variables:", err)
	// }

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
	}
	router.HandleRequests()
}
