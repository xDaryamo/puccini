package common

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var Quiet bool

func Fail(message string) {
	if !Quiet {
		fmt.Fprintln(color.Error, color.RedString(message))
	}
	os.Exit(1)
}

func Failf(f string, args ...interface{}) {
	Fail(fmt.Sprintf(f, args...))
}

func FailOnError(err error) {
	if err != nil {
		Failf("%s", err)
	}
}
