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

func Parse(urlString string, quirks []string, inputs map[string]interface{}) (*normal.ServiceTemplate, *problems.Problems, error) {

	context := NewContext(quirks)

	url_, err := url.NewValidURL(urlString, nil)
	if err != nil {
		return nil, &context.Problems, err
	}

	parserLock.Lock()

	// Phase 1: Read
	if !context.ReadServiceTemplate(url_) || !context.Problems.Empty() {
		parserLock.Unlock()
		return nil, &context.Problems, errors.New("phase 1: read")
	}

	// Phase 2: Namespaces
	context.AddNamespaces()
	if !context.Problems.Empty() {
		parserLock.Unlock()
		return nil, &context.Problems, errors.New("phase 2.1: namespaces")
	}
	context.LookupNames()
	if !context.Problems.Empty() {
		parserLock.Unlock()
		return nil, &context.Problems, errors.New("phase 2.2: namespaces lookup")
	}

	// Phase 3: Hieararchies
	context.AddHierarchies()
	if !context.Problems.Empty() {
		parserLock.Unlock()
		return nil, &context.Problems, errors.New("phase 3: hierarchies")
	}

	// Phase 4: Inheritance
	tasks := context.GetInheritTasks()
	tasks.Drain()
	if !context.Problems.Empty() {
		parserLock.Unlock()
		return nil, &context.Problems, errors.New("phase 4: inheritance")
	}

	SetInputs(context.ServiceTemplate.EntityPtr, inputs)

	// Phase 5: Rendering
	context.Render()
	if !context.Problems.Empty() {
		parserLock.Unlock()
		return nil, &context.Problems, errors.New("phase 5: rendering")
	}

	parserLock.Unlock()

	// Normalize
	s, ok := Normalize(context.ServiceTemplate.EntityPtr)
	if !ok || !context.Problems.Empty() {
		return nil, &context.Problems, errors.New("normalization")
	}

	return s, &context.Problems, nil
}
