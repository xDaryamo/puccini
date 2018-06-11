package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/tosca/compiler"
)

var output string

func init() {
	rootCmd.AddCommand(compileCmd)
	compileCmd.Flags().StringArrayVarP(&inputs, "input", "i", []string{}, "specify an input (name=JSON)")
	compileCmd.Flags().StringVarP(&output, "output", "o", "", "output Clout to file (instead of stdout)")
}

var compileCmd = &cobra.Command{
	Use:   "compile [[TOSCA PATH or URL]]",
	Short: "Compile TOSCA to Clout",
	Long:  `Parses a TOSCA service template and compiles the normalized output of the parser to Clout. Supports JavaScript plugins.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var urlString string
		if len(args) == 1 {
			urlString = args[0]
		}

		Compile(urlString)
	},
}

func Compile(urlString string) {
	s := Parse(urlString)

	c, err := compiler.Compile(s)
	common.ValidateError(err)

	if !common.Quiet || (output != "") {
		err = format.WriteOrPrint(c, ardFormat, true, output)
		common.ValidateError(err)
	}
}
