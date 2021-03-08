package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/ard"
	formatpkg "github.com/tliron/kutil/format"
	problemspkg "github.com/tliron/kutil/problems"
	"github.com/tliron/kutil/terminal"
	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/parser"
	"github.com/tliron/yamlkeys"
)

var inputs map[string]string
var inputsUrl string
var stopAtPhase uint32
var dumpPhases []uint
var filter string

var inputValues = make(map[string]interface{})

func init() {
	rootCommand.AddCommand(parseCommand)
	parseCommand.Flags().StringToStringVarP(&inputs, "input", "i", nil, "specify input (format is name=value)")
	parseCommand.Flags().StringVarP(&inputsUrl, "inputs", "n", "", "load inputs from a PATH or URL to YAML content")
	parseCommand.Flags().Uint32VarP(&stopAtPhase, "stop", "s", 5, "parser phase at which to end")
	parseCommand.Flags().UintSliceVarP(&dumpPhases, "dump", "d", nil, "dump phase internals")
	parseCommand.Flags().StringVarP(&filter, "filter", "t", "", "filter output by entity path; use '*' for wildcard matching (disables --stop and --dump)")
}

var parseCommand = &cobra.Command{
	Use:   "parse [[TOSCA PATH or URL]]",
	Short: "Parse TOSCA",
	Long:  `Parses and validates a TOSCA service template and reports problems if there are any. Provides access to phase diagnostics and entities.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var url string
		if len(args) == 1 {
			url = args[0]
		}

		if filter != "" {
			stopAtPhase = 5
			dumpPhases = nil
		}

		_, s := Parse(url)

		if (filter == "") && (len(dumpPhases) == 0) {
			formatpkg.Print(s, format, terminal.Stdout, strict, pretty)
		}
	},
}

func Parse(url string) (*parser.Context, *normal.ServiceTemplate) {
	ParseInputs()

	urlContext := urlpkg.NewContext()
	defer urlContext.Release()

	var url_ urlpkg.URL
	var err error
	if url == "" {
		log.Infof("parsing stdin", url)
		url_, err = urlpkg.ReadToInternalURLFromStdin("yaml")
	} else {
		log.Infof("parsing %q", url)
		url_, err = urlpkg.NewValidURL(url, nil, urlContext)
	}
	util.FailOnError(err)

	context := parser.NewContext(tosca.NewQuirks(quirks...))

	var problems *problemspkg.Problems

	// Phase 1: Read
	if stopAtPhase >= 1 {
		ok := context.ReadRoot(url_)

		context.MergeProblems()
		problems = context.GetProblems()
		FailOnProblems(problems)

		if !ok {
			// Stop here if failed to read
			util.Exit(1)
		}

		if ToPrintPhase(1) {
			if len(dumpPhases) > 1 {
				fmt.Fprintf(terminal.Stdout, "%s\n", terminal.StyleHeading("Imports"))
				context.PrintImports(1)
			} else {
				context.PrintImports(0)
			}
		}
	}

	// Phase 2: Namespaces
	if stopAtPhase >= 2 {
		context.AddNamespaces()
		context.LookupNames()
		if ToPrintPhase(2) {
			if len(dumpPhases) > 1 {
				fmt.Fprintf(terminal.Stdout, "%s\n", terminal.StyleHeading("Namespaces"))
				context.PrintNamespaces(1)
			} else {
				context.PrintNamespaces(0)
			}
		}
	}

	// Phase 3: Hieararchies
	if stopAtPhase >= 3 {
		context.AddHierarchies()
		if ToPrintPhase(3) {
			if len(dumpPhases) > 1 {
				fmt.Fprintf(terminal.Stdout, "%s\n", terminal.StyleHeading("Hierarchies"))
				context.PrintHierarchies(1)
			} else {
				context.PrintHierarchies(0)
			}
		}
	}

	// Phase 4: Inheritance
	if stopAtPhase >= 4 {
		tasks := context.GetInheritTasks()
		if ToPrintPhase(4) {
			if len(dumpPhases) > 1 {
				fmt.Fprintf(terminal.Stdout, "%s\n", terminal.StyleHeading("Inheritance Tasks"))
				tasks.Print(1)
			} else {
				tasks.Print(0)
			}
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
				fmt.Fprintf(terminal.Stdout, "%s\n", terminal.StyleHeading("Rendering"))
			}
			for _, entityPtr := range entityPtrs {
				fmt.Fprintf(terminal.Stdout, "%s:\n", terminal.StylePath(tosca.GetContext(entityPtr).Path.String()))
				err = formatpkg.Print(entityPtr, format, terminal.Stdout, strict, pretty)
				util.FailOnError(err)
			}
		}
	}

	if filter != "" {
		entityPtrs := context.Gather(filter)
		if len(entityPtrs) == 0 {
			util.Failf("No paths found matching filter: %q\n", filter)
		} else if !terminal.Quiet {
			for _, entityPtr := range entityPtrs {
				fmt.Fprintf(terminal.Stdout, "%s\n", terminal.StylePath(tosca.GetContext(entityPtr).Path.String()))
				err = formatpkg.Print(entityPtr, format, terminal.Stdout, strict, pretty)
				util.FailOnError(err)
			}
		}
	}

	context.MergeProblems()
	FailOnProblems(problems)

	// Normalize
	if serviceTemplate, ok := normal.NormalizeServiceTemplate(context.Root.EntityPtr); ok {
		return context, serviceTemplate
	} else {
		util.Fail("grammar does not support normalization")
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
		log.Infof("load inputs from %q", inputsUrl)

		urlContext := urlpkg.NewContext()
		defer urlContext.Release()

		url, err := urlpkg.NewValidURL(inputsUrl, nil, urlContext)
		util.FailOnError(err)
		reader, err := url.Open()
		util.FailOnError(err)
		defer reader.Close()
		data, err := yamlkeys.DecodeAll(reader)
		util.FailOnError(err)
		for _, data_ := range data {
			if map_, ok := data_.(ard.Map); ok {
				for key, value := range map_ {
					inputValues[yamlkeys.KeyString(key)] = value
				}
			} else {
				util.Failf("malformed inputs in %q", inputsUrl)
			}
		}
	}

	for name, input := range inputs {
		value, err := yamlkeys.Decode(strings.NewReader(input))
		util.FailOnError(err)
		inputValues[name] = value
	}
}
