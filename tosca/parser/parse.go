package parser

import (
	"errors"

	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/problems"
	urlpkg "github.com/tliron/puccini/url"
)

func Parse(url string, quirks tosca.Quirks, inputs map[string]interface{}) (*normal.ServiceTemplate, *problems.Problems, error) {
	context := NewContext(quirks)

	url_, err := urlpkg.NewValidURL(url, nil)
	if err != nil {
		return nil, nil, err
	}

	// Phase 1: Read
	ok := context.ReadRoot(url_)

	problems := context.GetProblems()

	if !problems.Empty() {
		return nil, problems, errors.New("phase 1: read")
	}

	if !ok {
		return nil, nil, errors.New("phase 1: read")
	}

	// Phase 2: Namespaces
	context.AddNamespaces()
	if !problems.Empty() {
		return nil, problems, errors.New("phase 2.1: namespaces")
	}
	context.LookupNames()
	if !problems.Empty() {
		return nil, problems, errors.New("phase 2.2: namespaces lookup")
	}

	// Phase 3: Hierarchies
	context.AddHierarchies()
	if !problems.Empty() {
		return nil, problems, errors.New("phase 3: hierarchies")
	}

	// Phase 4: Inheritance
	tasks := context.GetInheritTasks()
	tasks.Drain()
	if !problems.Empty() {
		return nil, problems, errors.New("phase 4: inheritance")
	}

	SetInputs(context.Root.EntityPtr, inputs)

	// Phase 5: Rendering
	context.Render()
	if !problems.Empty() {
		return nil, problems, errors.New("phase 5: rendering")
	}

	// Normalize
	serviceTemplate, ok := normal.NormalizeServiceTemplate(context.Root.EntityPtr)
	if !ok || !problems.Empty() {
		return nil, problems, errors.New("normalization")
	}

	return serviceTemplate, problems, nil
}
