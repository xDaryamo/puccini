package commands

import (
	"os"
	"runtime/pprof"

	"github.com/spf13/cobra"
	"github.com/tliron/commonlog"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

var logTo string
var verbose int
var cpuProfilePath string

func init() {
	rootCommand.PersistentFlags().BoolVarP(&terminal.Quiet, "quiet", "q", false, "suppress output")
	rootCommand.PersistentFlags().StringVarP(&logTo, "log", "l", "", "log to file (defaults to stderr)")
	rootCommand.PersistentFlags().CountVarP(&verbose, "verbose", "v", "add a log verbosity level (can be used twice)")
	rootCommand.PersistentFlags().StringVarP(&cpuProfilePath, "cpu-profile", "", "", "CPU profile file path")
}

var rootCommand = &cobra.Command{
	Use:   toolName,
	Short: "CSAR packaging tool",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cleanup, err := terminal.ProcessColorizeFlag(colorize)
		util.FailOnError(err)
		if cleanup != nil {
			util.OnExitError(cleanup)
		}

		if logTo == "" {
			if terminal.Quiet {
				verbose = -4
			}
			commonlog.Configure(verbose, nil)
		} else {
			commonlog.Configure(verbose, &logTo)
		}

		if cpuProfilePath != "" {
			cpuProfile, err := os.Create(cpuProfilePath)
			util.FailOnError(err)
			err = pprof.StartCPUProfile(cpuProfile)
			util.FailOnError(err)
			util.OnExit(pprof.StopCPUProfile)
		}
	},
}

func Execute() {
	util.FailOnError(rootCommand.Execute())
}
