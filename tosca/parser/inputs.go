package parser

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
)

func (self *ServiceContext) SetInputs(inputs map[string]ard.Value) {
	self.Context.entitiesLock.Lock()
	defer self.Context.entitiesLock.Unlock()

	tosca.SetInputs(self.Root.EntityPtr, inputs)
}
