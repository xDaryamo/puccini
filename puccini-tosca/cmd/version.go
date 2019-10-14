package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/puccini-tosca/version"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of puccini-tosca",
	Long:  `Shows the version of puccini-tosca.`,
	Run: func(cmd *cobra.Command, args []string) {
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			fmt.Fprintf(format.Stdout, "module.version=%s\n", buildInfo.Main.Version)
		}
		if version.GitRevision != "" {
			fmt.Fprintf(format.Stdout, "git.revision=%s\n", version.GitRevision)
		}
	},
}
