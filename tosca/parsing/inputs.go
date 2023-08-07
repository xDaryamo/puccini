package parsing

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/reflection"
)

//
// HasInputs
//

type HasInputs interface {
	SetInputs(map[string]ard.Value)
}

// From HasInputs interface
func SetInputs(entityPtr EntityPtr, inputs map[string]ard.Value) bool {
	if inputs == nil {
		return false
	}

	var done bool

	reflection.TraverseEntities(entityPtr, false, func(entityPtr EntityPtr) bool {
		if hasInputs, ok := entityPtr.(HasInputs); ok {
			hasInputs.SetInputs(inputs)
			done = true

			// Only one entity should implement the interface
			return false
		}
		return true
	})

	return done
}
