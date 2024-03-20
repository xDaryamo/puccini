package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(scriptletCommand)
}

var scriptletCommand = &cobra.Command{
	Use:   "scriptlet",
	Short: "Manage and process JavaScript scriptlets for Clout",
}
