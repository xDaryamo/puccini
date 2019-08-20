package cmd

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/parser"
	"github.com/tliron/puccini/url"
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
		var urlString string
		if len(args) == 1 {
			urlString = args[0]
		}

		if filter != "" {
			stopAtPhase = 5
			dumpPhases = nil
		}

		_, s := Parse(urlString)

		if (filter == "") && (len(dumpPhases) == 0) {
			format.Print(s, ardFormat, pretty)
		}
	},
}

func Parse(urlString string) (parser.Context, *normal.ServiceTemplate) {
	ParseInputs()

	var url_ url.URL
	var err error
	if urlString == "" {
		log.Infof("parsing stdin", urlString)
		url_, err = url.ReadInternalURLFromStdin("yaml")
	} else {
		log.Infof("parsing %s", urlString)
		url_, err = url.NewValidURL(urlString, nil)
	}
	common.FailOnError(err)

	context := parser.NewContext(quirks)

	// Phase 1: Read
	if stopAtPhase >= 1 {
		if !context.ReadServiceTemplate(url_) {
			// Stop here if there are errors
			if !common.Quiet {
				context.Problems.Print()
			}
			os.Exit(1)
		}

		// Stop here if there are problems
		if !context.Problems.Empty() {
			if !common.Quiet {
				context.Problems.Print()
			}
			os.Exit(1)
		}

		if !common.Quiet && ToPrintPhase(1) {
			if len(dumpPhases) > 1 {
				fmt.Fprintf(format.Stdout, "%s\n", format.ColorHeading("Imports"))
			}
			context.PrintImports(1)
		}
	}

	// Phase 2: Namespaces
	if stopAtPhase >= 2 {
		context.AddNamespaces()
		context.LookupNames()
		if !common.Quiet && ToPrintPhase(2) {
			if len(dumpPhases) > 1 {
				fmt.Fprintf(format.Stdout, "%s\n", format.ColorHeading("Namespaces"))
			}
			context.PrintNamespaces(1)
		}
	}

	// Phase 3: Hieararchies
	if stopAtPhase >= 3 {
		context.AddHierarchies()
		if !common.Quiet && ToPrintPhase(3) {
			if len(dumpPhases) > 1 {
				fmt.Fprintf(format.Stdout, "%s\n", format.ColorHeading("Hierarchies"))
			}
			context.PrintHierarchies(1)
		}
	}

	// Phase 4: Inheritance
	if stopAtPhase >= 4 {
		tasks := context.GetInheritTasks()
		if !common.Quiet && ToPrintPhase(4) {
			if len(dumpPhases) > 1 {
				fmt.Fprintf(format.Stdout, "%s\n", format.ColorHeading("Inheritance Tasks"))
			}
			tasks.Print(1)
		}
		tasks.Drain()
	}

	parser.SetInputs(context.ServiceTemplate.EntityPtr, inputValues)

	// Phase 5: Rendering
	if stopAtPhase >= 5 {
		entityPtrs := context.Render()
		if !common.Quiet && ToPrintPhase(5) {
			sort.Sort(entityPtrs)
			if len(dumpPhases) > 1 {
				fmt.Fprintf(format.Stdout, "%s\n", format.ColorHeading("Rendering"))
			}
			for _, entityPtr := range entityPtrs {
				fmt.Fprintf(format.Stdout, "%s:\n", format.ColorPath(tosca.GetContext(entityPtr).Path.String()))
				err = format.Print(entityPtr, ardFormat, pretty)
				common.FailOnError(err)
			}
		}
	}

	if filter != "" {
		entityPtrs := context.Gather(filter)
		if len(entityPtrs) == 0 {
			common.Failf("No paths found matching filter: \"%s\"\n", filter)
		} else if !common.Quiet {
			for _, entityPtr := range entityPtrs {
				fmt.Fprintf(format.Stdout, "%s\n", format.ColorPath(tosca.GetContext(entityPtr).Path.String()))
				err = format.Print(entityPtr, ardFormat, pretty)
				common.FailOnError(err)
			}
		}
	}

	if !common.Quiet {
		context.Problems.Print()
	}

	if !context.Problems.Empty() {
		os.Exit(1)
	}

	// Normalize
	s, ok := parser.Normalize(context.ServiceTemplate.EntityPtr)
	if !ok {
		common.Fail("grammar does not support normalization")
	}

	return context, s
}

func ToPrintPhase(phase uint) bool {
	for _, p := range dumpPhases {
		if p == phase {
			return true
		}
	}
	return false
}

func ParseInputs() {
	if inputsUrl != "" {
		log.Infof("load inputs from %s", inputsUrl)
		url_, err := url.NewValidURL(inputsUrl, nil)
		common.FailOnError(err)
		reader, err := url_.Open()
		common.FailOnError(err)
		if readerCloser, ok := reader.(io.ReadCloser); ok {
			defer readerCloser.Close()
		}
		data, err := format.Read(reader, "yaml")
		common.FailOnError(err)
		if map_, ok := data.(ard.Map); ok {
			for key, value := range map_ {
				inputValues[key] = value
			}
		} else {
			common.Failf("malformed inputs in %s", inputsUrl)
		}
	}

	for _, input := range inputs {
		s := strings.SplitN(input, "=", 2)
		if len(s) != 2 {
			common.Failf("malformed input: %s", input)
		}
		value, err := format.Decode(s[1], "yaml")
		common.FailOnError(err)
		inputValues[s[0]] = value
	}
}
