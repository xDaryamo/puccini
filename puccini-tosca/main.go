package main

import (
	"github.com/tebeka/atexit"
	"github.com/tliron/puccini/puccini-tosca/commands"
)

var BuildCommit string

func main() {
	commands.Execute()
	atexit.Exit(0)
}
