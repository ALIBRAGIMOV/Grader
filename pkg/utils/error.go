package utils

import "log"

func FatalOnError(msg string, err error) {
	if err != nil {
		log.Fatal(msg + " " + err.Error())
	}
}
