package common

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var Quiet bool

func Error(message string) {
	if !Quiet {
		fmt.Fprintln(color.Error, color.RedString(message))
	}
	os.Exit(1)
}

func Errorf(f string, args ...interface{}) {
	Error(fmt.Sprintf(f, args...))
}

func ValidateError(err error) {
	if err != nil {
		Errorf("%s", err)
	}
}
