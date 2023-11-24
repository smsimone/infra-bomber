package utils

import "log"

func HandleErrFail(err error) {
	if err != nil {
		log.Fatalf("Got error: %v", err)
	}
}
