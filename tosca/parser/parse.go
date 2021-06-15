package parser

import (
	"errors"
	"sync"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/problems"
	"github.com/tliron/kutil/terminal"
	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

var parserLock sync.Mutex

func Parse(url urlpkg.URL, stylist *terminal.Stylist, quirks tosca.Quirks, inputs map[string]ard.Value) (*Context, *normal.ServiceTemplate, *problems.Problems, error) {
	context := NewContext(stylist, quirks)

	// Phase 1: Read
	parserLock.Lock()
	ok := context.ReadRoot(url, "")
	parserLock.Unlock()

	context.MergeProblems()
	problems := context.GetProblems()

	if !problems.Empty() {
		return context, nil, problems, errors.New("read problems")
	}

	if !ok {
		return context, nil, nil, errors.New("read error")
	}

	// Phase 2: Namespaces
	parserLock.Lock()
	context.AddNamespaces()
	parserLock.Unlock()
	parserLock.Lock()
	context.LookupNames()
	parserLock.Unlock()

	// Phase 3: Hierarchies
	parserLock.Lock()
	context.AddHierarchies()
	parserLock.Unlock()

	// Phase 4: Inheritance
	parserLock.Lock()
	tasks := context.GetInheritTasks()
	tasks.Drain()
	parserLock.Unlock()

	SetInputs(context.Root.EntityPtr, inputs)

	// Phase 5: Rendering
	parserLock.Lock()
	context.Render()
	parserLock.Unlock()

	context.MergeProblems()
	if !problems.Empty() {
		return context, nil, problems, errors.New("parsing problems")
	}

	// Normalize
	parserLock.Lock()
	serviceTemplate, ok := normal.NormalizeServiceTemplate(context.Root.EntityPtr)
	parserLock.Unlock()
	if !ok || !problems.Empty() {
		return context, nil, problems, errors.New("normalization")
	}

	return context, serviceTemplate, problems, nil
}
