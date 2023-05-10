package utils

import (
	"log"
)

func CheckErr(err error, fbs ...func()) {
	var onError, onSuccess func()
	if len(fbs) == 1 {
		onError = fbs[0]
	} else if len(fbs) == 2 {
		onError = fbs[0]
		onSuccess = fbs[1]
	}
	if err != nil {
		log.Printf("Error: %v", err)
		if onError != nil {
			onError()
		}
	} else {
		if onSuccess != nil {
			onSuccess()
		}
	}
}
