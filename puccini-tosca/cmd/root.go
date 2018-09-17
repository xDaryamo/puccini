package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/format"
)

var logTo string
var verbosity int
var ardFormat string
var quirks []string

var bashCompletionTo string

func init() {
	rootCmd.PersistentFlags().BoolVarP(&common.Quiet, "quiet", "q", false, "suppress output")
	rootCmd.PersistentFlags().StringVarP(&logTo, "log", "l", "", "log to file (defaults to stderr)")
	rootCmd.PersistentFlags().CountVarP(&verbosity, "verbosity", "v", "add a log verbosity level (can be used twice)")
	rootCmd.PersistentFlags().StringVarP(&ardFormat, "format", "f", "", "force format (\"yaml\", \"json\", or \"xml\")")
	rootCmd.PersistentFlags().StringSliceVarP(&quirks, "quirk", "x", nil, "parser quirk")

	rootCmd.Flags().StringVarP(&bashCompletionTo, "bash-completion", "b", "", "generate bash completion file")
}

var rootCmd = &cobra.Command{
	Use:   "puccini-tosca",
	Short: "TOSCA frontend for Clout",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if logTo == "" {
			if common.Quiet {
				verbosity = -4
			}
			common.ConfigureLogging(verbosity, nil)
		} else {
			common.ConfigureLogging(verbosity, &logTo)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if bashCompletionTo != "" {
			if !common.Quiet {
				fmt.Fprintf(format.Stdout, "generating bash completion script: %s\n", bashCompletionTo)
			}
			cmd.GenBashCompletionFile(bashCompletionTo)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	common.ValidateError(err)
}
