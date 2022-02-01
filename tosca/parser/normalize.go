package parser

import (
	"github.com/tliron/puccini/tosca/normal"
)

func (self *ServiceContext) Normalize() (*normal.ServiceTemplate, bool) {
	self.Context.entitiesLock.Lock()
	defer self.Context.entitiesLock.Unlock()

	return normal.NormalizeServiceTemplate(self.Root.EntityPtr)
}
