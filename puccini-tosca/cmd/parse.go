package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/compiler"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/parser"
	"github.com/tliron/puccini/url"
)

var inputs []string
var stopAtPhase uint32
var printPhases []uint
var coerce bool
var examine string

var inputValues = make(map[string]interface{})

func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.Flags().StringArrayVarP(&inputs, "input", "i", []string{}, "specify an input (name=value)")
	parseCmd.Flags().Uint32VarP(&stopAtPhase, "stop", "s", 6, "parser phase at which to stop")
	parseCmd.Flags().UintSliceVarP(&printPhases, "print", "p", nil, "parser phases to print")
	parseCmd.Flags().BoolVarP(&coerce, "coerce", "c", false, "emit final values (calls intrinsic functions)")
	parseCmd.Flags().StringVarP(&examine, "examine", "e", "", "examine entities with path, may use '*' for wildcards (disables --print)")
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

		if examine != "" {
			// Examine cancels printing phases
			printPhases = nil
		}

		s := Parse(urlString)

		if (examine == "") && (len(printPhases) == 0) {
			format.Print(s, ardFormat, true)
		}
	},
}

func Parse(urlString string) *normal.ServiceTemplate {
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
	common.ValidateError(err)

	var s *normal.ServiceTemplate

	context := parser.NewContext()

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
			if len(printPhases) > 1 {
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
			if len(printPhases) > 1 {
				fmt.Fprintf(format.Stdout, "%s\n", format.ColorHeading("Namespaces"))
			}
			context.PrintNamespaces(1)
		}
	}

	// Phase 3: Hieararchies
	if stopAtPhase >= 3 {
		context.AddHierarchies()
		if !common.Quiet && ToPrintPhase(3) {
			if len(printPhases) > 1 {
				fmt.Fprintf(format.Stdout, "%s\n", format.ColorHeading("Hierarchies"))
			}
			context.PrintHierarchies(1)
		}
	}

	// Phase 4: Inheritance
	if stopAtPhase >= 4 {
		tasks := context.GetInheritTasks()
		if !common.Quiet && ToPrintPhase(4) {
			if len(printPhases) > 1 {
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
			if len(printPhases) > 1 {
				fmt.Fprintf(format.Stdout, "%s\n", format.ColorHeading("Rendering"))
			}
			for _, entityPtr := range entityPtrs {
				fmt.Fprintf(format.Stdout, "%s:\n", format.ColorPath(tosca.GetContext(entityPtr).Path))
				err = format.Print(entityPtr, ardFormat, true)
				common.ValidateError(err)
			}
		}
	}

	// Phase 6: Topology
	if stopAtPhase >= 6 {
		var ok bool
		if s, ok = parser.Normalize(context.ServiceTemplate.EntityPtr); ok {
			if coerce {
				// Check for coercion problems
				clout, err := compiler.Compile(s)
				common.ValidateError(err)
				compiler.Coerce(clout, context.Problems)
			}

			// Only print if there are no problems
			if !common.Quiet && ToPrintPhase(6) && context.Problems.Empty() {
				fmt.Fprintf(format.Stdout, "%s\n", format.ColorHeading("Topology"))
				s.PrintRelationships(1)
			}
		}
	}

	if examine != "" {
		entityPtrs := context.Gather(examine)
		if len(entityPtrs) == 0 {
			common.Errorf("Examine path not found: \"%s\"\n", examine)
		} else if !common.Quiet {
			for _, entityPtr := range entityPtrs {
				if len(entityPtrs) > 0 {
					fmt.Fprintf(format.Stdout, "%s\n", format.ColorPath(tosca.GetContext(entityPtr).Path))
				}
				err = format.Print(entityPtr, ardFormat, true)
				common.ValidateError(err)
			}
		}
	}

	if !common.Quiet {
		context.Problems.Print()
	}

	if !context.Problems.Empty() {
		os.Exit(1)
	}

	return s
}

func ToPrintPhase(phase uint) bool {
	for _, p := range printPhases {
		if p == phase {
			return true
		}
	}
	return false
}

func ParseInputs() {
	for _, input := range inputs {
		s := strings.SplitN(input, "=", 2)
		if len(s) != 2 {
			common.Errorf("malformed input: %s", input)
		}
		value, err := format.Decode(s[1], "yaml")
		common.ValidateError(err)
		inputValues[s[0]] = value
	}
}
