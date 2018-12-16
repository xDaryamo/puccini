package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/js"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list [[Clout PATH or URL]]",
	Short: "List JavaScript in Clout",
	Long:  ``,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		if len(args) == 1 {
			path = args[0]
		}

		clout_, err := ReadClout(path)
		common.FailOnError(err)

		List(clout_)
	},
}

func List(clout_ *clout.Clout) {
	metadata, err := js.GetMetadata(clout_)
	common.FailOnError(err)

	ListValue(metadata, nil)
}

func ListValue(value interface{}, path []string) {
	switch v := value.(type) {
	case string:
		if !common.Quiet {
			fmt.Fprintf(format.Stdout, "%s\n", strings.Join(path, "."))
		}
	case ard.Map:
		for key, vv := range v {
			ListValue(vv, append(path, key))
		}
	}
}
