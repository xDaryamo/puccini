package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tliron/puccini/ard"
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/clout/js"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/common/terminal"
)

func init() {
	rootCommand.AddCommand(listCommand)
}

var listCommand = &cobra.Command{
	Use:   "list [[Clout PATH or URL]]",
	Short: "List JavaScript scriptlets in Clout",
	Long:  ``,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		if len(args) == 1 {
			path = args[0]
		}

		clout, err := ReadClout(path)
		common.FailOnError(err)

		List(clout)
	},
}

func List(clout *cloutpkg.Clout) {
	metadata, err := js.GetScriptletsMetadata(clout)
	common.FailOnError(err)

	ListValue(metadata, nil)
}

func ListValue(value interface{}, path []string) {
	switch value_ := value.(type) {
	case string:
		if !terminal.Quiet {
			fmt.Fprintf(terminal.Stdout, "%s\n", strings.Join(path, "."))
		}

	case ard.StringMap:
		for key, value__ := range value_ {
			ListValue(value__, append(path, key))
		}
	}
}
