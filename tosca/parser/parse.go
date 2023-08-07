package parser

import (
	contextpkg "context"
	"errors"

	"github.com/tliron/exturl"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/problems"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// ParseContext
//

type ParseContext struct {
	URL     exturl.URL
	Origins []exturl.URL
	Stylist *terminal.Stylist
	Quirks  parsing.Quirks
	Inputs  map[string]ard.Value
}

//
// Result
//

type Result struct {
	ServiceContext        *ServiceContext
	NormalServiceTemplate *normal.ServiceTemplate
	Problems              *problems.Problems
}

func (self *Context) Parse(context contextpkg.Context, parseContext ParseContext) (Result, error) {
	var result Result

	result.ServiceContext = self.NewServiceContext(parseContext.Stylist, parseContext.Quirks)

	// Phase 1: Read
	ok := result.ServiceContext.ReadRoot(context, parseContext.URL, parseContext.Origins, "")
	result.ServiceContext.MergeProblems()
	result.Problems = result.ServiceContext.GetProblems()

	if !result.Problems.Empty() {
		return result, errors.New("read problems")
	}

	if !ok {
		return result, errors.New("read error")
	}

	// Phase 2: Namespaces
	result.ServiceContext.AddNamespaces()
	result.ServiceContext.LookupNames()

	// Phase 3: Hierarchies
	result.ServiceContext.AddHierarchies()

	// Phase 4: Inheritance
	result.ServiceContext.Inherit(nil)

	result.ServiceContext.SetInputs(parseContext.Inputs)

	// Phase 5: Rendering
	result.ServiceContext.Render()
	result.ServiceContext.MergeProblems()
	if !result.Problems.Empty() {
		return result, errors.New("parsing problems")
	}

	// Normalize
	result.NormalServiceTemplate, ok = result.ServiceContext.Normalize()
	if !ok || !result.Problems.Empty() {
		return result, errors.New("normalization")
	}

	return result, nil
}
