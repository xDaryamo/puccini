package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tebeka/atexit"
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/common"
	formatpkg "github.com/tliron/puccini/common/format"
	problemspkg "github.com/tliron/puccini/common/problems"
	"github.com/tliron/puccini/common/terminal"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/parser"
	urlpkg "github.com/tliron/puccini/url"
	"github.com/tliron/yamlkeys"
)

var inputs []string
var inputsUrl string
var stopAtPhase uint32
var dumpPhases []uint
var filter string

var inputValues = make(map[string]interface{})

func init() {
	rootCommand.AddCommand(parseCommand)
	parseCommand.Flags().StringArrayVarP(&inputs, "input", "i", []string{}, "specify an input (name=YAML)")
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

func Parse(url string) (parser.Context, *normal.ServiceTemplate) {
	ParseInputs()

	var url_ urlpkg.URL
	var err error
	if url == "" {
		log.Infof("parsing stdin", url)
		url_, err = urlpkg.ReadToInternalURLFromStdin("yaml")
	} else {
		log.Infof("parsing \"%s\"", url)
		url_, err = urlpkg.NewValidURL(url, nil)
	}
	common.FailOnError(err)
	defer url_.Release()

	context := parser.NewContext(tosca.NewQuirks(quirks...))
	defer context.Release()
	var problems *problemspkg.Problems

	// Phase 1: Read
	if stopAtPhase >= 1 {
		ok := context.ReadRoot(url_)

		problems = context.GetProblems()
		FailOnProblems(problems)

		if !ok {
			// Stop here if failed to read
			atexit.Exit(1)
		}

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
				err = formatpkg.Print(entityPtr, format, terminal.Stdout, strict, pretty)
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
				err = formatpkg.Print(entityPtr, format, terminal.Stdout, strict, pretty)
				common.FailOnError(err)
			}
		}
	}

	FailOnProblems(problems)

	// Normalize
	if serviceTemplate, ok := normal.NormalizeServiceTemplate(context.Root.EntityPtr); ok {
		return context, serviceTemplate
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
		url, err := urlpkg.NewValidURL(inputsUrl, nil)
		common.FailOnError(err)
		defer url.Release()
		reader, err := url.Open()
		common.FailOnError(err)
		defer reader.Close()
		data, err := formatpkg.ReadYAML(reader)
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
		value, err := formatpkg.DecodeYAML(s[1])
		common.FailOnError(err)
		inputValues[s[0]] = value
	}
}
