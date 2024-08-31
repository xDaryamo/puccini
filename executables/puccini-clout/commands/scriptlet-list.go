package commands

import (
	contextpkg "context"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tliron/exturl"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/terminal"
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

		urlContext := exturl.NewContext()
		util.OnExitError(urlContext.Release)

		context, cancel := contextpkg.WithTimeout(contextpkg.Background(), time.Duration(timeout*float64(time.Second)))
		util.OnExit(cancel)

		clout := LoadClout(context, url, urlContext)

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
