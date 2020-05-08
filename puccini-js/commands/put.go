package commands

import (
	"github.com/spf13/cobra"
	"github.com/tliron/puccini/clout/js"
	"github.com/tliron/puccini/common"
	formatpkg "github.com/tliron/puccini/common/format"
	"github.com/tliron/puccini/common/terminal"
	urlpkg "github.com/tliron/puccini/url"
)

func init() {
	rootCommand.AddCommand(putCommand)
	putCommand.Flags().StringVarP(&output, "output", "o", "", "output Clout to file (default is stdout)")
}

var putCommand = &cobra.Command{
	Use:   "put [NAME] [JavaScript PATH or URL] [[Clout PATH or URL]]",
	Short: "Put JavaScript scriptlet in Clout",
	Long:  ``,
	Args:  cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		scriptletName := args[0]
		jsUrl := args[1]

		var cloutPath string
		if len(args) == 3 {
			cloutPath = args[2]
		}

		clout, err := ReadClout(cloutPath)
		common.FailOnError(err)

		url, err := urlpkg.NewValidURL(jsUrl, nil)
		common.FailOnError(err)
		defer url.Release()

		scriptlet, err := urlpkg.ReadToString(url)
		common.FailOnError(err)

		err = js.SetScriptlet(scriptletName, js.CleanupScriptlet(scriptlet), clout)
		common.FailOnError(err)

		err = formatpkg.WriteOrPrint(clout, format, terminal.Stdout, strict, pretty, output)
		common.FailOnError(err)
	},
}
