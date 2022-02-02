package parser

import (
	"github.com/tliron/puccini/tosca/normal"
)

func (self *ServiceContext) Normalize() (*normal.ServiceTemplate, bool) {
	self.Context.lock.Lock()
	defer self.Context.lock.Unlock()

	return normal.NormalizeServiceTemplate(self.Root.EntityPtr)
}
