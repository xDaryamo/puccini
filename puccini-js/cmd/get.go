package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/js"
)

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&output, "output", "o", "", "output to file (instead of stdout)")
}

var getCmd = &cobra.Command{
	Use:   "get [COMMAND] [[Clout PATH or URL]]",
	Short: "Get JavaScript from Clout",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		var path string
		if len(args) == 2 {
			path = args[1]
		}

		clout, err := ReadClout(path)
		common.FailOnError(err)

		sourceCode, err := js.GetScriptSourceCode(name, clout)
		common.FailOnError(err)

		if !common.Quiet {
			err = format.WriteOrPrint(sourceCode, ardFormat, pretty, output)
			common.FailOnError(err)
		}
	},
}
