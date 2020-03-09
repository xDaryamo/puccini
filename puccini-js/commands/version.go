package commands

import (
	"github.com/tliron/puccini/version"
)

func init() {
	rootCommand.AddCommand(version.NewCommand("puccini-js"))
}
