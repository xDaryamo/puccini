package common

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var Quiet bool

func Errorf(f string, args ...interface{}) {
	if !Quiet {
		fmt.Fprintln(color.Error, color.RedString(fmt.Sprintf(f, args...)))
	}
	os.Exit(1)
}

func ValidateError(err error) {
	if err != nil {
		Errorf("%s", err)
	}
}
