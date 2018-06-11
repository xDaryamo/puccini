package common

import (
	"fmt"
	"os"
)

var Quiet bool

func Errorf(message string, args ...interface{}) {
	if !Quiet {
		fmt.Fprintf(os.Stderr, message+"\n", args...)
	}
	os.Exit(1)
}

func ValidateError(err error) {
	if err != nil {
		Errorf("%s", err)
	}
}
