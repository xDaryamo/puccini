package parser

import (
	"errors"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/common/problems"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	urlpkg "github.com/tliron/puccini/url"
)

func Parse(url urlpkg.URL, quirks tosca.Quirks, inputs map[string]ard.Value) (*normal.ServiceTemplate, *problems.Problems, error) {
	context := NewContext(quirks)

	// Phase 1: Read
	ok := context.ReadRoot(url)

	context.MergeProblems()
	problems := context.GetProblems()

	if !problems.Empty() {
		return nil, problems, errors.New("read problems")
	}

	if !ok {
		return nil, nil, errors.New("read error")
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
		return nil, problems, errors.New("parsing problems")
	}

	// Normalize
	serviceTemplate, ok := normal.NormalizeServiceTemplate(context.Root.EntityPtr)
	if !ok || !problems.Empty() {
		return nil, problems, errors.New("normalization")
	}

	return serviceTemplate, problems, nil
}
