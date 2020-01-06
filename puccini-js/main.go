package main

import (
	"github.com/tebeka/atexit"
	"github.com/tliron/puccini/puccini-js/cmd"
)

func main() {
	cmd.Execute()
	atexit.Exit(0)
}
