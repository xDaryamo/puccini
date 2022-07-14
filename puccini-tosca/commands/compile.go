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

var output string
var resolve bool
var coerce bool
var exec string
var arguments map[string]string

func init() {
	rootCommand.AddCommand(compileCommand)
	compileCommand.Flags().StringSliceVarP(&importPaths, "path", "b", nil, "specify an import path or base URL")
	compileCommand.Flags().StringVarP(&template, "template", "t", "", "select service template in CSAR (leave empty for root, or use \"all\", path, or integer index)")
	compileCommand.Flags().StringToStringVarP(&inputs, "input", "i", nil, "specify input (format is name=value)")
	compileCommand.Flags().StringVarP(&inputsUrl, "inputs", "n", "", "load inputs from a PATH or URL to YAML content")
	compileCommand.Flags().StringVarP(&output, "output", "o", "", "output Clout to file (leave empty for stdout)")
	compileCommand.Flags().BoolVarP(&resolve, "resolve", "r", true, "resolves the topology (attempts to satisfy all requirements with capabilities)")
	compileCommand.Flags().BoolVarP(&coerce, "coerce", "c", false, "coerces all values (calls functions and applies constraints)")
	compileCommand.Flags().StringVarP(&exec, "exec", "e", "", "execute JavaScript scriptlet")
	compileCommand.Flags().StringToStringVarP(&arguments, "argument", "a", nil, "used with --exec to specify a scriptlet argument (format is key=value)")
}

var compileCommand = &cobra.Command{
	Use:   "compile [[TOSCA PATH or URL]]",
	Short: "Compile TOSCA to Clout",
	Long:  `Parses TOSCA service templates and compiles the normalized output to Clout.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var url string
		if len(args) == 1 {
			url = args[0]
		}

		Compile(url)
	},
}

func Compile(url string) {
	// Parse
	context, serviceTemplate := Parse(url)
	problems := context.GetProblems()
	urlContext := context.Root.GetContext().URL.Context()

	// Compile
	clout, err := serviceTemplate.Compile()
	util.FailOnError(err)

	// Resolve
	if resolve {
		js.Resolve(clout, problems, urlContext, true, format, strict, pretty)
		FailOnProblems(problems)
	}

	// Coerce
	if coerce {
		js.Coerce(clout, problems, urlContext, true, format, strict, pretty)
		FailOnProblems(problems)
	}

	if exec != "" {
		err = Exec(exec, arguments, clout, urlContext)
		util.FailOnError(err)
	} else if !terminal.Quiet || (output != "") {
		err = transcribe.WriteOrPrint(clout, format, terminal.Stdout, strict, pretty, output)
		util.FailOnError(err)
	}
}

func Exec(scriptletName string, arguments map[string]string, clout *cloutpkg.Clout, urlContext *urlpkg.Context) error {
	// Try loading JavaScript from Clout
	scriptlet, err := js.GetScriptlet(scriptletName, clout)

	if err != nil {
		urlContext := urlpkg.NewContext()
		defer urlContext.Release()

		// Try loading JavaScript from path or URL
		url, err := urlpkg.NewValidURL(scriptletName, nil, urlContext)
		util.FailOnError(err)

		scriptlet, err = urlpkg.ReadString(url)
		util.FailOnError(err)

		err = js.SetScriptlet(exec, js.CleanupScriptlet(scriptlet), clout)
		util.FailOnError(err)
	}

	jsContext := js.NewContext(scriptletName, log, arguments, terminal.Quiet, format, strict, pretty, output, urlContext)
	_, err = jsContext.Require(clout, scriptletName, nil)
	return js.UnwrapException(err)
}
