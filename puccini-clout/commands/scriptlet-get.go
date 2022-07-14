package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/transcribe"
	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/kutil/util"
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/clout/js"
)

func init() {
	scriptletCommand.AddCommand(getCommand)
	getCommand.Flags().StringVarP(&output, "output", "o", "", "output to file (default is stdout)")
}

var getCommand = &cobra.Command{
	Use:   "get [NAME] [[Clout PATH or URL]]",
	Short: "Get JavaScript scriptlet from Clout",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		scriptletName := args[0]

		var url string
		if len(args) == 2 {
			url = args[1]
		}

		urlContext := urlpkg.NewContext()
		defer urlContext.Release()

		clout, err := cloutpkg.Load(url, inputFormat, urlContext)
		util.FailOnError(err)

		scriptlet, err := js.GetScriptlet(scriptletName, clout)
		util.FailOnError(err)

		if !terminal.Quiet {
			err = transcribe.WriteOrPrint(scriptlet, format, terminal.Stdout, strict, pretty, output)
			util.FailOnError(err)
		}
	},
}
