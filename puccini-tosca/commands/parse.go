package commands

import (
	"sort"

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

var importPaths []string
var template string
var inputs map[string]string
var inputsUrl string
var stopAtPhase uint32
var dumpPhases []uint
var filter string

var inputValues = make(map[string]any)

func init() {
	rootCommand.AddCommand(parseCommand)
	parseCommand.Flags().StringSliceVarP(&importPaths, "path", "b", nil, "specify an import path or base URL")
	parseCommand.Flags().StringVarP(&template, "template", "t", "", "select service template in CSAR (leave empty for root, or use path or integer index)")
	parseCommand.Flags().StringToStringVarP(&inputs, "input", "i", nil, "specify an input (format is name=YAML)")
	parseCommand.Flags().StringVarP(&inputsUrl, "inputs", "n", "", "load inputs from a PATH or URL to YAML content")
	parseCommand.Flags().Uint32VarP(&stopAtPhase, "stop", "s", 5, "parser phase at which to end")
	parseCommand.Flags().UintSliceVarP(&dumpPhases, "dump", "d", nil, "dump phase internals")
	parseCommand.Flags().StringVarP(&filter, "filter", "r", "", "filter output by entity path; use '*' for wildcard matching (disables --stop and --dump)")
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

		_, serviceTemplate := Parse(url)

		if (filter == "") && (len(dumpPhases) == 0) {
			formatpkg.Print(serviceTemplate, format, terminal.Stdout, strict, pretty)
		}
	},
}

var parserContext = parser.NewContext()

func Parse(url string) (*parser.ServiceContext, *normal.ServiceTemplate) {
	ParseInputs()

	urlContext := urlpkg.NewContext()
	util.OnExitError(urlContext.Release)

	var origins []urlpkg.URL
	for _, importPath := range importPaths {
		origin, err := urlpkg.NewURL(importPath, urlContext)
		if err != nil {
			origin = urlpkg.NewFileURL(importPath, urlContext)
		}
		origins = append(origins, origin)
	}

	var url_ urlpkg.URL
	var err error
	if url == "" {
		log.Infof("parsing stdin", url)
		url_, err = urlpkg.ReadToInternalURLFromStdin("yaml")
	} else {
		log.Infof("parsing %q", url)
		url_, err = urlpkg.NewValidURL(url, origins, urlContext)
	}
	util.FailOnError(err)

	stylist := terminal.Stylize
	if problemsFormat != "" {
		stylist = terminal.NewStylist(false)
	}

	serviceContext := parserContext.NewServiceContext(stylist, tosca.NewQuirks(quirks...))

	var problems *problemspkg.Problems

	// Phase 1: Read
	if stopAtPhase >= 1 {
		ok := serviceContext.ReadRoot(url_, origins, template)

		serviceContext.MergeProblems()
		problems = serviceContext.GetProblems()
		FailOnProblems(problems)

		if !ok {
			// Stop here if failed to read
			util.Exit(1)
		}

		if ToPrintPhase(1) {
			if len(dumpPhases) > 1 {
				terminal.Printf("%s\n", terminal.Stylize.Heading("Imports"))
				serviceContext.PrintImports(1)
			} else {
				serviceContext.PrintImports(0)
			}
		}
	}

	// Phase 2: Namespaces
	if stopAtPhase >= 2 {
		serviceContext.AddNamespaces()
		serviceContext.LookupNames()
		if ToPrintPhase(2) {
			if len(dumpPhases) > 1 {
				terminal.Printf("%s\n", terminal.Stylize.Heading("Namespaces"))
				serviceContext.PrintNamespaces(1)
			} else {
				serviceContext.PrintNamespaces(0)
			}
		}
	}

	// Phase 3: Hieararchies
	if stopAtPhase >= 3 {
		serviceContext.AddHierarchies()
		if ToPrintPhase(3) {
			if len(dumpPhases) > 1 {
				terminal.Printf("%s\n", terminal.Stylize.Heading("Hierarchies"))
				serviceContext.PrintHierarchies(1)
			} else {
				serviceContext.PrintHierarchies(0)
			}
		}
	}

	// Phase 4: Inheritance
	if stopAtPhase >= 4 {
		if ToPrintPhase(4) {
			serviceContext.Inherit(func(tasks parser.Tasks) {
				if len(dumpPhases) > 1 {
					terminal.Printf("%s\n", terminal.Stylize.Heading("Inheritance Tasks"))
					tasks.Print(1)
				} else {
					tasks.Print(0)
				}
			})
		} else {
			serviceContext.Inherit(nil)
		}
	}

	if serviceContext.Root == nil {
		return serviceContext, nil
	}

	serviceContext.SetInputs(inputValues)

	// Phase 5: Rendering
	if stopAtPhase >= 5 {
		entityPtrs := serviceContext.Render()
		if ToPrintPhase(5) {
			sort.Sort(entityPtrs)
			if len(dumpPhases) > 1 {
				terminal.Printf("%s\n", terminal.Stylize.Heading("Rendering"))
			}
			for _, entityPtr := range entityPtrs {
				terminal.Printf("%s:\n", terminal.Stylize.Path(tosca.GetContext(entityPtr).Path.String()))
				err = formatpkg.Print(entityPtr, format, terminal.Stdout, strict, pretty)
				util.FailOnError(err)
			}
		}
	}

	if filter != "" {
		entityPtrs := serviceContext.Gather(filter)
		if len(entityPtrs) == 0 {
			util.Failf("No paths found matching filter: %q\n", filter)
		} else if !terminal.Quiet {
			for _, entityPtr := range entityPtrs {
				terminal.Printf("%s\n", terminal.Stylize.Path(tosca.GetContext(entityPtr).Path.String()))
				err = formatpkg.Print(entityPtr, format, terminal.Stdout, strict, pretty)
				util.FailOnError(err)
			}
		}
	}

	serviceContext.MergeProblems()
	FailOnProblems(problems)

	// Normalize
	if serviceTemplate, ok := serviceContext.Normalize(); ok {
		return serviceContext, serviceTemplate
	} else {
		util.Fail("grammar does not support normalization")
		return serviceContext, nil
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
		util.OnExitError(urlContext.Release)

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

	if inputs != nil {
		for name, input := range inputs {
			input_, _, err := ard.DecodeYAML(input, false)
			util.FailOnError(err)
			inputValues[name] = input_
		}
	}
}
