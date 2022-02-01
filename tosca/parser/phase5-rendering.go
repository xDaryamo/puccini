package parser

import (
	"github.com/tliron/puccini/tosca"
)

func (self *ServiceContext) Render() tosca.EntityPtrs {
	self.Context.entitiesLock.Lock()
	defer self.Context.entitiesLock.Unlock()

	var entityPtrs tosca.EntityPtrs

	self.Context.renderWork.TraverseEntities(self.Root.EntityPtr, func(entityPtr tosca.EntityPtr) bool {
		if tosca.Render(entityPtr) {
			entityPtrs = append(entityPtrs, entityPtr)
		}
		return true
	})

	return entityPtrs
}
