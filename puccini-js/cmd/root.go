package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/format"
)

var logTo string
var verbose int
var ardFormat string
var pretty bool

var bashCompletionTo string

func init() {
	rootCmd.PersistentFlags().BoolVarP(&common.Quiet, "quiet", "q", false, "suppress output")
	rootCmd.PersistentFlags().StringVarP(&logTo, "log", "l", "", "log to file (defaults to stderr)")
	rootCmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "add a log verbosity level (can be used twice)")
	rootCmd.PersistentFlags().StringVarP(&ardFormat, "format", "f", "", "force format (\"yaml\", \"json\", or \"xml\")")
	rootCmd.PersistentFlags().BoolVarP(&pretty, "pretty", "p", true, "prettify output")

	rootCmd.Flags().StringVarP(&bashCompletionTo, "bash-completion", "b", "", "generate bash completion file")
}

var rootCmd = &cobra.Command{
	Use:   "puccini-js",
	Short: "JavaScript processor for Clout",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if logTo == "" {
			if common.Quiet {
				verbose = -4
			}
			common.ConfigureLogging(verbose, nil)
		} else {
			common.ConfigureLogging(verbose, &logTo)
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
	common.FailOnError(err)
}
