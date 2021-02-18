package parser

import (
	"errors"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/problems"
	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

func Parse(url urlpkg.URL, quirks tosca.Quirks, inputs map[string]ard.Value) (*Context, *normal.ServiceTemplate, *problems.Problems, error) {
	context := NewContext(quirks)

	// Phase 1: Read
	ok := context.ReadRoot(url)

	context.MergeProblems()
	problems := context.GetProblems()

	if !problems.Empty() {
		return context, nil, problems, errors.New("read problems")
	}

	if !ok {
		return context, nil, nil, errors.New("read error")
	}

	// Phase 2: Namespaces
	context.AddNamespaces()
	context.LookupNames()

	// Phase 3: Hierarchies
	context.AddHierarchies()

	// Phase 4: Inheritance
	tasks := context.GetInheritTasks()
	tasks.Drain()

	SetInputs(context.Root.EntityPtr, inputs)

	// Phase 5: Rendering
	context.Render()

	context.MergeProblems()
	if !problems.Empty() {
		return context, nil, problems, errors.New("parsing problems")
	}

	// Normalize
	serviceTemplate, ok := normal.NormalizeServiceTemplate(context.Root.EntityPtr)
	if !ok || !problems.Empty() {
		return context, nil, problems, errors.New("normalization")
	}

	return context, serviceTemplate, problems, nil
}
