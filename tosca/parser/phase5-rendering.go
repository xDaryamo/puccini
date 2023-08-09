package parser

import (
	"github.com/tliron/puccini/tosca/parsing"
)

func (self *Context) Render() parsing.EntityPtrs {
	self.Parser.lock.Lock()
	defer self.Parser.lock.Unlock()

	var entityPtrs parsing.EntityPtrs

	self.Parser.renderWork.TraverseEntities(self.Root.EntityPtr, func(entityPtr parsing.EntityPtr) bool {
		if parsing.Render(entityPtr) {
			entityPtrs = append(entityPtrs, entityPtr)
		}
		return true
	})

	return entityPtrs
}
