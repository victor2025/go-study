package utils

import (
	"log"
)

// err: error object
// cbs: callback functions, [0]:onError [1]:onSuccess
func HandleError(err error, cbs ...func()) {
	var onError, onSuccess func()
	if len(cbs) == 1 {
		onError = cbs[0]
	} else if len(cbs) == 2 {
		onError = cbs[0]
		onSuccess = cbs[1]
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
