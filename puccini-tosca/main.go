package main

import (
	"github.com/tebeka/atexit"
	"github.com/tliron/puccini/puccini-tosca/cmd"
)

var BuildCommit string

func main() {
	cmd.Execute()
	atexit.Exit(0)
}
