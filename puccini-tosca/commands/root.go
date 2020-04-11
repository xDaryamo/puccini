package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/common/terminal"
)

var logTo string
var verbose int
var format string
var colorize bool
var strict bool
var timestamps bool
var pretty bool
var problemsFormat string
var quirks []string

func init() {
	rootCommand.PersistentFlags().BoolVarP(&terminal.Quiet, "quiet", "q", false, "suppress output")
	rootCommand.PersistentFlags().StringVarP(&logTo, "log", "l", "", "log to file (defaults to stderr)")
	rootCommand.PersistentFlags().CountVarP(&verbose, "verbose", "v", "add a log verbosity level (can be used twice)")
	rootCommand.PersistentFlags().StringVarP(&format, "format", "f", "", "force output format (\"yaml\", \"json\", or \"xml\")")
	rootCommand.PersistentFlags().BoolVarP(&colorize, "colorize", "z", true, "colorize output")
	rootCommand.PersistentFlags().BoolVarP(&strict, "strict", "y", false, "strict output (for \"YAML\" format only)")
	rootCommand.PersistentFlags().BoolVarP(&timestamps, "timestamps", "w", true, "allow timestamps (for \"YAML\" format only)")
	rootCommand.PersistentFlags().BoolVarP(&pretty, "pretty", "p", true, "prettify output")
	rootCommand.PersistentFlags().StringVarP(&problemsFormat, "problems-format", "m", "", "problems format (\"yaml\", \"json\", or \"xml\")")
	rootCommand.PersistentFlags().StringSliceVarP(&quirks, "quirk", "x", nil, "parser quirk")
}

var rootCommand = &cobra.Command{
	Use:   toolName,
	Short: "TOSCA frontend for Clout",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if colorize {
			terminal.EnableColor()
		}
		if logTo == "" {
			if terminal.Quiet {
				verbose = -4
			}
			common.ConfigureLogging(verbose, nil)
		} else {
			common.ConfigureLogging(verbose, &logTo)
		}
	},
}

func Execute() {
	err := rootCommand.Execute()
	common.FailOnError(err)
}
