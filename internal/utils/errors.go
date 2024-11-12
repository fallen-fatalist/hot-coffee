package utils

import (
	"fmt"
	"os"
)

func FatalError(context string, err error) {
	if err != nil {
		fmt.Printf(context+": %s\n", err)
		os.Exit(1)
	}
}
