package helpers

import (
	"log"
)

func FailedOnError(err error, message string) {
	log.Panicf("%s: %s", err.Error(), message)
}
