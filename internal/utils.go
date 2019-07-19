package internal

import (
	"log"
	"runtime/debug"
)

func check(err error) bool {
	if err != nil {
		debug.PrintStack()
		log.Println(err)
		return true
	}
	return false
}
