package parser

import (
	"github.com/tliron/puccini/normal"
)

func (self *Context) Normalize() (*normal.ServiceTemplate, bool) {
	self.Parser.lock.Lock()
	defer self.Parser.lock.Unlock()

	return normal.NormalizeServiceTemplate(self.Root.EntityPtr)
}
