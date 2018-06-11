package parser

import (
	"github.com/tliron/puccini/tosca/reflection"
)

//
// HasInputs
//

type HasInputs interface {
	SetInputs(map[string]interface{})
}

// From HasInputs interface
func SetInputs(entityPtr interface{}, inputs map[string]interface{}) {
	reflection.Traverse(entityPtr, func(entityPtr interface{}) bool {
		if hasInputs, ok := entityPtr.(HasInputs); ok {
			hasInputs.SetInputs(inputs)

			// Only one entity should implement the interface
			return false
		}
		return true
	})
}
