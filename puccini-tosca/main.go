package main

import (
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/puccini-tosca/commands"

	_ "github.com/tliron/kutil/logging/simple"
)

func main() {
	commands.Execute()
	util.Exit(0)
}
