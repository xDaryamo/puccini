package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/commonlog"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

var (
	logTo          string
	verbose        int
	format         string
	colorize       string
	strict         bool
	pretty         bool
	base64         bool
	timeout        float64
	cpuProfilePath string
)

func init() {
	rootCommand.PersistentFlags().BoolVarP(&terminal.Quiet, "quiet", "q", false, "suppress output")
	rootCommand.PersistentFlags().StringVarP(&logTo, "log", "l", "", "log to file (defaults to stderr)")
	rootCommand.PersistentFlags().CountVarP(&verbose, "verbose", "v", "add a log verbosity level (can be used twice)")
	rootCommand.PersistentFlags().BoolVarP(&commonlog.Trace, "trace", "", false, "add stack trace to log messages")
	rootCommand.PersistentFlags().StringVarP(&format, "format", "f", "", "force output format (\"yaml\", \"json\", \"xjson\", \"xml\", \"cbor\", \"messagepack\", or \"go\")")
	rootCommand.PersistentFlags().StringVarP(&colorize, "colorize", "z", "true", "colorize output (boolean or \"force\")")
	rootCommand.PersistentFlags().BoolVarP(&strict, "strict", "y", false, "strict output (for \"yaml\" format only)")
	rootCommand.PersistentFlags().BoolVarP(&pretty, "pretty", "p", true, "prettify output")
	rootCommand.PersistentFlags().BoolVarP(&base64, "base64", "", false, "output base64 (for \"cbor\", \"messagepack\" formats)")
	rootCommand.PersistentFlags().Float64Var(&timeout, "timeout", 30.0, "timeout in seconds")
	rootCommand.PersistentFlags().StringVarP(&cpuProfilePath, "cpu-profile", "", "", "CPU profile file path")
}

var rootCommand = &cobra.Command{
	Use:   toolName,
	Short: "TOSCA parser and compiler",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		util.InitializeCPUProfiling(cpuProfilePath)
		util.InitializeColorization(colorize)
		commonlog.Initialize(verbose, logTo)
	},
}

func Execute() {
	util.FailOnError(rootCommand.Execute())
}
