package parser

import (
	"fmt"

	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/problems"
	"github.com/tliron/puccini/url"
)

func Parse(urlString string, inputs map[string]interface{}) (*normal.ServiceTemplate, *problems.Problems, error) {
	context := NewContext()

	url_, err := url.NewValidURL(urlString, nil)
	if err != nil {
		return nil, context.Problems, err
	}

	// Phase 1: Read
	if !context.ReadServiceTemplate(url_) || !context.Problems.Empty() {
		return nil, context.Problems, fmt.Errorf("phase 1")
	}

	// Phase 2: Namespaces
	context.AddNamespaces()
	if !context.Problems.Empty() {
		return nil, context.Problems, fmt.Errorf("phase 2 (namespaces)")
	}
	context.LookupNames()
	if !context.Problems.Empty() {
		return nil, context.Problems, fmt.Errorf("phase 2 (lookups)")
	}

	// Phase 3: Hieararchies
	context.AddHierarchies()
	if !context.Problems.Empty() {
		return nil, context.Problems, fmt.Errorf("phase 3")
	}

	// Phase 4: Inheritance
	tasks := context.GetInheritTasks()
	tasks.Drain()
	if !context.Problems.Empty() {
		return nil, context.Problems, fmt.Errorf("phase 4")
	}

	SetInputs(context.ServiceTemplate.EntityPtr, inputs)

	// Phase 5: Rendering
	context.Render()
	if !context.Problems.Empty() {
		return nil, context.Problems, fmt.Errorf("phase 5")
	}

	// Phase 6: Topology
	s, ok := Normalize(context.ServiceTemplate.EntityPtr)
	if !ok || !context.Problems.Empty() {
		return nil, context.Problems, fmt.Errorf("phase 6")
	}

	return s, context.Problems, nil
}
