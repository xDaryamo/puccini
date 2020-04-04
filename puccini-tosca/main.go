package main

import (
	"github.com/tebeka/atexit"
	"github.com/tliron/puccini/puccini-tosca/commands"
)

func main() {
	commands.Execute()
	atexit.Exit(0)
}
