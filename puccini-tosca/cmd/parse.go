package cmd

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tebeka/atexit"
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/common"
	format_ "github.com/tliron/puccini/common/format"
	"github.com/tliron/puccini/common/terminal"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/parser"
	"github.com/tliron/puccini/tosca/problems"
	"github.com/tliron/puccini/url"
	"github.com/tliron/yamlkeys"
)

var inputs []string
var inputsUrl string
var stopAtPhase uint32
var dumpPhases []uint
var filter string

var inputValues = make(map[string]interface{})

func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.Flags().StringArrayVarP(&inputs, "input", "i", []string{}, "specify an input (name=YAML)")
	parseCmd.Flags().StringVarP(&inputsUrl, "inputs", "n", "", "load inputs from a PATH or URL to YAML content")
	parseCmd.Flags().Uint32VarP(&stopAtPhase, "stop", "s", 5, "parser phase at which to end")
	parseCmd.Flags().UintSliceVarP(&dumpPhases, "dump", "d", nil, "dump phase internals")
	parseCmd.Flags().StringVarP(&filter, "filter", "t", "", "filter output by entity path; use '*' for wildcard matching (disables --stop and --dump)")
}

var parseCmd = &cobra.Command{
	Use:   "parse [[TOSCA PATH or URL]]",
	Short: "Parse TOSCA",
	Long:  `Parses and validates a TOSCA service template and reports problems if there are any. Provides access to phase diagnostics and entities.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var url_ string
		if len(args) == 1 {
			url_ = args[0]
		}

		if filter != "" {
			stopAtPhase = 5
			dumpPhases = nil
		}

		_, s := Parse(url_)

		if (filter == "") && (len(dumpPhases) == 0) {
			format_.Print(s, format, pretty)
		}
	},
}

func Parse(url_ string) (parser.Context, *normal.ServiceTemplate) {
	ParseInputs()

	var url__ url.URL
	var err error
	if url_ == "" {
		log.Infof("parsing stdin", url_)
		url__, err = url.ReadToInternalURLFromStdin(format)
	} else {
		log.Infof("parsing \"%s\"", url_)
		url__, err = url.NewValidURL(url_, nil)
	}
	common.FailOnError(err)

	context := parser.NewContext(quirks)
	var problems_ *problems.Problems

	// Phase 1: Read
	if stopAtPhase >= 1 {
		if !context.ReadRoot(url__) {
			// Stop here if failed to read
			atexit.Exit(1)
		}

		problems_ = context.GetProblems()
		FailOnProblems(problems_)

		if ToPrintPhase(1) {
			if len(dumpPhases) > 1 {
				fmt.Fprintf(terminal.Stdout, "%s\n", terminal.ColorHeading("Imports"))
			}
			context.PrintImports(1)
		}
	}

	// Phase 2: Namespaces
	if stopAtPhase >= 2 {
		context.AddNamespaces()
		context.LookupNames()
		if ToPrintPhase(2) {
			if len(dumpPhases) > 1 {
				fmt.Fprintf(terminal.Stdout, "%s\n", terminal.ColorHeading("Namespaces"))
			}
			context.PrintNamespaces(1)
		}
	}

	// Phase 3: Hieararchies
	if stopAtPhase >= 3 {
		context.AddHierarchies()
		if ToPrintPhase(3) {
			if len(dumpPhases) > 1 {
				fmt.Fprintf(terminal.Stdout, "%s\n", terminal.ColorHeading("Hierarchies"))
			}
			context.PrintHierarchies(1)
		}
	}

	// Phase 4: Inheritance
	if stopAtPhase >= 4 {
		tasks := context.GetInheritTasks()
		if ToPrintPhase(4) {
			if len(dumpPhases) > 1 {
				fmt.Fprintf(terminal.Stdout, "%s\n", terminal.ColorHeading("Inheritance Tasks"))
			}
			tasks.Print(1)
		}
		tasks.Drain()
	}

	if context.Root == nil {
		return context, nil
	}

	parser.SetInputs(context.Root.EntityPtr, inputValues)

	// Phase 5: Rendering
	if stopAtPhase >= 5 {
		entityPtrs := context.Render()
		if ToPrintPhase(5) {
			sort.Sort(entityPtrs)
			if len(dumpPhases) > 1 {
				fmt.Fprintf(terminal.Stdout, "%s\n", terminal.ColorHeading("Rendering"))
			}
			for _, entityPtr := range entityPtrs {
				fmt.Fprintf(terminal.Stdout, "%s:\n", terminal.ColorPath(tosca.GetContext(entityPtr).Path.String()))
				err = format_.Print(entityPtr, format, pretty)
				common.FailOnError(err)
			}
		}
	}

	if filter != "" {
		entityPtrs := context.Gather(filter)
		if len(entityPtrs) == 0 {
			common.Failf("No paths found matching filter: \"%s\"\n", filter)
		} else if !terminal.Quiet {
			for _, entityPtr := range entityPtrs {
				fmt.Fprintf(terminal.Stdout, "%s\n", terminal.ColorPath(tosca.GetContext(entityPtr).Path.String()))
				err = format_.Print(entityPtr, format, pretty)
				common.FailOnError(err)
			}
		}
	}

	FailOnProblems(problems_)

	// Normalize
	if s, ok := normal.NormalizeServiceTemplate(context.Root.EntityPtr); ok {
		return context, s
	} else {
		common.Fail("grammar does not support normalization")
		return context, nil
	}
}

func ToPrintPhase(phase uint) bool {
	if !terminal.Quiet {
		for _, phase_ := range dumpPhases {
			if phase_ == phase {
				return true
			}
		}
	}
	return false
}

func ParseInputs() {
	if inputsUrl != "" {
		log.Infof("load inputs from \"%s\"", inputsUrl)
		url_, err := url.NewValidURL(inputsUrl, nil)
		common.FailOnError(err)
		reader, err := url_.Open()
		common.FailOnError(err)
		if readerCloser, ok := reader.(io.ReadCloser); ok {
			defer readerCloser.Close()
		}
		data, err := format_.Read(reader, "yaml")
		common.FailOnError(err)
		if map_, ok := data.(ard.Map); ok {
			for key, value := range map_ {
				inputValues[yamlkeys.KeyString(key)] = value
			}
		} else {
			common.Failf("malformed inputs in \"%s\"", inputsUrl)
		}
	}

	for _, input := range inputs {
		s := strings.SplitN(input, "=", 2)
		if len(s) != 2 {
			common.Failf("malformed input: %s", input)
		}
		value, err := format_.Decode(s[1], "yaml")
		common.FailOnError(err)
		inputValues[s[0]] = value
	}
}
