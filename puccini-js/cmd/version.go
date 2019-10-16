package cmd

import (
	"github.com/tliron/puccini/version"
)

func init() {
	rootCmd.AddCommand(version.NewCommand("puccini-js"))
}
