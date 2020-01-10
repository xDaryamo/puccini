package parser

import (
	"errors"
	"sync"

	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/problems"
	"github.com/tliron/puccini/url"
)

// TODO: this is a brute-force method to avoid race conditions
var parserLock sync.Mutex

func Parse(url_ string, quirks []string, inputs map[string]interface{}) (*normal.ServiceTemplate, *problems.Problems, error) {
	context := NewContext(quirks)

	url__, err := url.NewValidURL(url_, nil)
	if err != nil {
		return nil, nil, err
	}

	parserLock.Lock()
	defer parserLock.Unlock()

	// Phase 1: Read
	if !context.ReadRoot(url__) {
		return nil, nil, errors.New("phase 1: read")
	}

	problems := context.GetProblems()

	if !problems.Empty() {
		return nil, problems, errors.New("phase 1: read")
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
	s, ok := normal.NormalizeServiceTemplate(context.Root.EntityPtr)
	if !ok || !problems.Empty() {
		return nil, problems, errors.New("normalization")
	}

	return s, problems, nil
}
