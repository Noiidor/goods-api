package errors

import "log"

var AsyncErrors = make(chan error)

func HandleAsyncErrors() {
	go func() {
		for {
			err := <-AsyncErrors
			if err == nil {
				continue
			}
			log.Printf("ASYNC ERROR: %v", err)
		}
	}()
}
