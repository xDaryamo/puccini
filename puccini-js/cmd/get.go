package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tliron/puccini/common"
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

		c, err := ReadClout(path)
		common.ValidateError(err)

		sourceCode, err := js.GetScriptSourceCode(name, c)
		common.ValidateError(err)

		if !common.Quiet {
			// TODO: write to file?
			fmt.Println(sourceCode)
		}
	},
}
