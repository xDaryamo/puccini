package common

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var Quiet bool

func Errorf(message string, args ...interface{}) {
	if !Quiet {
		fmt.Fprintf(color.Error, message+"\n", args...)
	}
	os.Exit(1)
}

func ValidateError(err error) {
	if err != nil {
		Errorf("%s", err)
	}
}
