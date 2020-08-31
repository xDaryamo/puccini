package parser

import (
	"github.com/tliron/kutil/reflection"
	"github.com/tliron/puccini/tosca"
)

func (self *Context) Render() tosca.EntityPtrs {
	return Render(self.Root.EntityPtr)
}

//
// Renderable
//

type Renderable interface {
	Render()
}

// From Renderable interface
func Render(entityPtr tosca.EntityPtr) tosca.EntityPtrs {
	var entityPtrs tosca.EntityPtrs

	reflection.Traverse(entityPtr, func(entityPtr tosca.EntityPtr) bool {
		if renderable, ok := entityPtr.(Renderable); ok {
			renderable.Render()
			entityPtrs = append(entityPtrs, entityPtr)
		}
		return true
	})

	return entityPtrs
}
