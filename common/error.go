package common

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/tebeka/atexit"
)

var Quiet bool

func Fail(message string) {
	if !Quiet {
		fmt.Fprintln(color.Error, color.RedString(message))
	}
	atexit.Exit(1)
}

func Failf(f string, args ...interface{}) {
	Fail(fmt.Sprintf(f, args...))
}

func FailOnError(err error) {
	if err != nil {
		Failf("%s", err)
	}
}
