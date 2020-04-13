package commands

import (
	"github.com/spf13/cobra"
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/clout/js"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/common/terminal"
	urlpkg "github.com/tliron/puccini/url"
)

func init() {
	rootCommand.AddCommand(execCommand)
	execCommand.Flags().StringVarP(&output, "output", "o", "", "output to file or directory (default is stdout)")
}

var execCommand = &cobra.Command{
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

		clout, err := ReadClout(path)
		common.FailOnError(err)

		// Try loading JavaScript from Clout
		scriptlet, err := js.GetScriptlet(scriptletName, clout)

		if err != nil {
			// Try loading JavaScript from path or URL
			url, err := urlpkg.NewValidURL(scriptletName, nil)
			common.FailOnError(err)

			scriptlet, err = urlpkg.ReadToString(url)
			common.FailOnError(err)

			err = js.SetScriptlet(scriptletName, js.CleanupScriptlet(scriptlet), clout)
			common.FailOnError(err)
		}

		err = Exec(scriptletName, scriptlet, clout)
		common.FailOnError(err)
	},
}

func Exec(scriptletName string, scriptlet string, clout *cloutpkg.Clout) error {
	jsContext := js.NewContext(scriptletName, log, terminal.Quiet, format, strict, timestamps, pretty, output)

	program, err := jsContext.GetProgram(scriptletName, scriptlet)
	if err != nil {
		return err
	}

	runtime := jsContext.NewCloutRuntime(clout, nil)

	_, err = runtime.RunProgram(program)

	return js.UnwrapException(err)
}
