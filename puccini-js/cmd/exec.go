package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/js"
)

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().StringVarP(&output, "output", "o", "", "output to file or directory (default is stdout)")
}

var execCmd = &cobra.Command{
	Use:   "exec [COMMAND] [[Clout PATH or URL]]",
	Short: "Execute JavaScript in Clout",
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

		err = Exec(name, sourceCode, c)
		common.ValidateError(err)
	},
}

func Exec(name string, sourceCode string, c *clout.Clout) error {
	program, err := js.GetProgram(name, sourceCode)
	if err != nil {
		return err
	}

	context := js.NewContext(name, log, common.Quiet, ardFormat, output)
	_, runtime := context.NewCloutContext(c)
	_, err = runtime.RunProgram(program)

	return js.UnwrapError(err)
}
