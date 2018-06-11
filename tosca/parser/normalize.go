package parser

import (
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/reflection"
)

//
// Normalizable
//

type Normalizable interface {
	Normalize() *normal.ServiceTemplate
}

// From Normalizable interface
func Normalize(entityPtr interface{}) (*normal.ServiceTemplate, bool) {
	var s *normal.ServiceTemplate

	reflection.Traverse(entityPtr, func(entityPtr interface{}) bool {
		if normalizable, ok := entityPtr.(Normalizable); ok {
			s = normalizable.Normalize()

			// Only one entity should implement the interface
			return false
		}
		return true
	})

	if s == nil {
		return nil, false
	}

	return s, true
}
