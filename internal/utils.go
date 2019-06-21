package internal

import "log"

func check(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}
