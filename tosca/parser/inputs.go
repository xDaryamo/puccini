package parser

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

func (self *ServiceContext) SetInputs(inputs map[string]ard.Value) {
	self.Context.lock.Lock()
	defer self.Context.lock.Unlock()

	parsing.SetInputs(self.Root.EntityPtr, inputs)
}
