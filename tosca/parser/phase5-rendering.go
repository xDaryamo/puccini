package parser

import (
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/reflection"
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
func Render(entityPtr interface{}) tosca.EntityPtrs {
	var entityPtrs tosca.EntityPtrs

	reflection.Traverse(entityPtr, func(entityPtr interface{}) bool {
		if renderable, ok := entityPtr.(Renderable); ok {
			lock := common.GetLock(entityPtr)
			lock.Lock()
			defer lock.Unlock()

			renderable.Render()
			entityPtrs = append(entityPtrs, entityPtr)
		}
		return true
	})

	return entityPtrs
}
