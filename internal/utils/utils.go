package utils

import (
	"log"
	"os"
)

func FatalError(context string, err error) {
	if err != nil {
		log.Fatalf(context+": %s\n", err)
		os.Exit(1)
	}
}

func In(str string, arr []string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
