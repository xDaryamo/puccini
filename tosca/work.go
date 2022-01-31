package tosca

import (
	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/reflection"
)

//
// EntityWork
//

type EntityWork map[EntityPtr]struct{}

func (self EntityWork) Start(log logging.Logger, entityPtr EntityPtr) bool {
	if _, ok := self[entityPtr]; ok {
		log.Debugf("skip: %s", GetContext(entityPtr).Path)
		return true
	}
	self[entityPtr] = struct{}{}
	return false
}

func (self EntityWork) TraverseEntities(entityPtr EntityPtr, traverse reflection.EntityTraverser) {
	reflection.TraverseEntities(entityPtr, false, func(entityPtr EntityPtr) bool {
		if self.Start(log, entityPtr) {
			return false
		}
		return traverse(entityPtr)
	})
}
