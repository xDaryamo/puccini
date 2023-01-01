package commands

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/transcribe"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca/csar"

	urlpkg "github.com/tliron/kutil/url"
)

var output string
var format string
var colorize string
var strict bool
var pretty bool

func init() {
	rootCommand.AddCommand(metaCommand)

	metaCommand.Flags().StringVarP(&archiveFormat, "archive-format", "a", "", "force archive format (\"tar.gz\", \"tar\", or \"zip\"); leave empty to determine automatically from extension")
	metaCommand.Flags().StringVarP(&output, "output", "o", "", "output metadata to file (leave empty for stdout)")
	metaCommand.Flags().StringVarP(&format, "format", "f", "", "force output format (\"yaml\", \"json\", \"cjson\", \"xml\", \"cbor\", \"messagepack\", or \"go\")")
	metaCommand.Flags().StringVarP(&colorize, "colorize", "z", "true", "colorize output (boolean or \"force\")")
	metaCommand.Flags().BoolVarP(&strict, "strict", "y", false, "strict output (for \"yaml\" format only)")
	metaCommand.Flags().BoolVarP(&pretty, "pretty", "p", true, "prettify output")
}

var metaCommand = &cobra.Command{
	Use:   "meta [CSAR PATH or URL]",
	Short: "Show CSAR metadata",
	Long:  `Parses, validates, and extracts CSAR metadata.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var url = args[0]

		Meta(url)
	},
}

func Meta(url string) {
	var err error

	urlContext := urlpkg.NewContext()
	util.OnExitError(urlContext.Release)

	var csarUrl urlpkg.URL
	csarUrl, err = urlpkg.NewValidURL(url, nil, urlContext)
	util.FailOnError(err)

	var meta *csar.Meta
	meta, err = csar.ReadMetaFromURL(csarUrl, archiveFormat)
	util.FailOnError(err)

	if !terminal.Quiet || (output != "") {
		err = transcribe.WriteOrPrint(meta, format, os.Stdout, strict, pretty, output)
		util.FailOnError(err)
	}
}
