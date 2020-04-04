package parser

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/common/reflection"
	"github.com/tliron/puccini/tosca"
)

//
// HasInputs
//

type HasInputs interface {
	SetInputs(map[string]ard.Value)
}

// From HasInputs interface
func SetInputs(entityPtr tosca.EntityPtr, inputs map[string]ard.Value) {
	if inputs == nil {
		return
	}

	reflection.Traverse(entityPtr, func(entityPtr tosca.EntityPtr) bool {
		if hasInputs, ok := entityPtr.(HasInputs); ok {
			hasInputs.SetInputs(inputs)

			// Only one entity should implement the interface
			return false
		}
		return true
	})
}
