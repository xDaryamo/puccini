package common

import (
	"fmt"

	"github.com/tebeka/atexit"
	"github.com/tliron/puccini/common/terminal"
)

func Fail(message string) {
	if !terminal.Quiet {
		fmt.Fprintln(terminal.Stderr, terminal.ColorError(message))
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
