package main

import (
	"git.csnzoo.com/mduguay/mockex/internal"
)

func main() {

	var sto internal.Storage
	sto.Connect()
	defer sto.Disconnect()

	router := internal.Router{
		Stg: sto,
	}

	router.HandleRequests()
}
