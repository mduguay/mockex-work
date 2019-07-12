package internal

import "log"

func check(err error) bool {
	if err != nil {
		log.Fatalln(err)
		return true
	}
	return false
}
