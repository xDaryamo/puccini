package parser

import (
	"github.com/tliron/puccini/tosca/parsing"
)

func (self *ServiceContext) Render() parsing.EntityPtrs {
	self.Context.lock.Lock()
	defer self.Context.lock.Unlock()

	var entityPtrs parsing.EntityPtrs

	self.Context.renderWork.TraverseEntities(self.Root.EntityPtr, func(entityPtr parsing.EntityPtr) bool {
		if parsing.Render(entityPtr) {
			entityPtrs = append(entityPtrs, entityPtr)
		}
		return true
	})

	return entityPtrs
}
