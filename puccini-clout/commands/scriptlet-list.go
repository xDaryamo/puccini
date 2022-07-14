package commands

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/terminal"
	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/kutil/util"
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/clout/js"
)

func init() {
	scriptletCommand.AddCommand(listCommand)
}

var listCommand = &cobra.Command{
	Use:   "list [[Clout PATH or URL]]",
	Short: "List JavaScript scriptlets in Clout",
	Long:  ``,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var url string
		if len(args) == 1 {
			url = args[0]
		}

		urlContext := urlpkg.NewContext()
		defer urlContext.Release()

		clout, err := cloutpkg.Load(url, inputFormat, urlContext)
		util.FailOnError(err)

		List(clout)
	},
}

func List(clout *cloutpkg.Clout) {
	metadata, err := js.GetScriptletsMetadata(clout)
	util.FailOnError(err)

	ListValue(metadata, nil)
}

func ListValue(value any, path []string) {
	switch value_ := value.(type) {
	case string:
		if !terminal.Quiet {
			terminal.Printf("%s\n", strings.Join(path, "."))
		}

	case ard.StringMap:
		for key, value__ := range value_ {
			ListValue(value__, append(path, key))
		}
	}
}
