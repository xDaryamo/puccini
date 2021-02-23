package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

var logTo string
var verbose int
var format string
var colorize string
var strict bool
var timestamps bool
var pretty bool
var problemsFormat string
var quirks []string

func init() {
	rootCommand.PersistentFlags().BoolVarP(&terminal.Quiet, "quiet", "q", false, "suppress output")
	rootCommand.PersistentFlags().StringVarP(&logTo, "log", "l", "", "log to file (defaults to stderr)")
	rootCommand.PersistentFlags().CountVarP(&verbose, "verbose", "v", "add a log verbosity level (can be used twice)")
	rootCommand.PersistentFlags().StringVarP(&format, "format", "f", "", "force output format (\"yaml\", \"json\", \"cjson\", \"xml\", or \"cbor\")")
	rootCommand.PersistentFlags().StringVarP(&colorize, "colorize", "z", "true", "colorize output (boolean or \"force\")")
	rootCommand.PersistentFlags().BoolVarP(&strict, "strict", "y", false, "strict output (for \"YAML\" format only)")
	rootCommand.PersistentFlags().BoolVarP(&timestamps, "timestamps", "w", false, "use !!timestamp type (for \"YAML\" format)")
	rootCommand.PersistentFlags().BoolVarP(&pretty, "pretty", "p", true, "prettify output")
	rootCommand.PersistentFlags().StringVarP(&problemsFormat, "problems-format", "m", "", "problems format (\"yaml\", \"json\", \"cjson\", \"xml\", or \"cbor\")")
	rootCommand.PersistentFlags().StringSliceVarP(&quirks, "quirk", "x", nil, "parser quirk")
}

var rootCommand = &cobra.Command{
	Use:   toolName,
	Short: "TOSCA frontend for Clout",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := terminal.ProcessColorizeFlag(colorize)
		util.FailOnError(err)
		if logTo == "" {
			if terminal.Quiet {
				verbose = -4
			}
			logging.Configure(verbose, nil)
		} else {
			logging.Configure(verbose, &logTo)
		}
	},
}

func Execute() {
	err := rootCommand.Execute()
	util.FailOnError(err)
}
