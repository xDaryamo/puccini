package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/url"
)

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().StringVarP(&output, "output", "o", "", "output to file or directory (default is stdout)")
}

var execCmd = &cobra.Command{
	Use:   "exec [NAME or JavaScript PATH or URL] [[Clout PATH or URL]]",
	Short: "Execute JavaScript scriptlet in Clout",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		scriptletName := args[0]

		var path string
		if len(args) == 2 {
			path = args[1]
		}

		clout_, err := ReadClout(path)
		common.FailOnError(err)

		// Try loading JavaScript from Clout
		scriptlet, err := js.GetScriptlet(scriptletName, clout_)

		if err != nil {
			// Try loading JavaScript from path or URL
			url_, err := url.NewValidURL(scriptletName, nil)
			common.FailOnError(err)

			scriptlet, err = url.Read(url_)
			common.FailOnError(err)

			err = js.SetScriptlet(scriptletName, js.CleanupScriptlet(scriptlet), clout_)
			common.FailOnError(err)
		}

		err = Exec(scriptletName, scriptlet, clout_)
		common.FailOnError(err)
	},
}

func Exec(scriptletName string, scriptlet string, clout_ *clout.Clout) error {
	jsContext := js.NewContext(scriptletName, log, common.Quiet, ardFormat, pretty, output)

	program, err := jsContext.GetProgram(scriptletName, scriptlet)
	if err != nil {
		return err
	}

	runtime := jsContext.NewCloutRuntime(clout_, nil)

	_, err = runtime.RunProgram(program)

	return js.UnwrapException(err)
}
