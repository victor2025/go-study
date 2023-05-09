package utils

import (
	"log"
)

func HandleError(err error, fbs ...func()) {
	var onError, onSuccess func()
	if len(fbs) == 1 {
		onError = fbs[0]
	} else if len(fbs) == 2 {
		onError = fbs[0]
		onSuccess = fbs[1]
	}
	if err != nil {
		log.Fatalf("Fatal: %v", err)
		if onError != nil {
			onError()
		}
	} else {
		if onSuccess != nil {
			onSuccess()
		}
	}
}
