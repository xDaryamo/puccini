package parser

import (
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

var renderWork = make(tosca.EntityWork)

// From Renderable interface
func Render(entityPtr tosca.EntityPtr) tosca.EntityPtrs {
	var entityPtrs tosca.EntityPtrs

	renderWork.TraverseEntities(entityPtr, func(entityPtr tosca.EntityPtr) bool {
		if renderable, ok := entityPtr.(Renderable); ok {
			renderable.Render()
			entityPtrs = append(entityPtrs, entityPtr)
		}
		return true
	})

	return entityPtrs
}
