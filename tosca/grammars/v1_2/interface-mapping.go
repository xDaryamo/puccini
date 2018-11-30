package v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/v1_1"
)

//
// InterfaceMapping
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.11
//

type InterfaceMapping struct {
	*v1_1.Entity `name:"interface mapping"`

	NodeTemplateName *string `require:"0"`
	InterfaceName    *string `require:"1"`

	NodeTemplate *v1_1.NodeTemplate `lookup:"0,NodeTemplateName" json:"-" yaml:"-"`
}

func NewInterfaceMapping(context *tosca.Context) *InterfaceMapping {
	return &InterfaceMapping{Entity: v1_1.NewEntity(context)}
}

// tosca.Reader signature
func ReadInterfaceMapping(context *tosca.Context) interface{} {
	self := NewInterfaceMapping(context)
	if context.ValidateType("list") {
		strings := context.ReadStringListFixed(2)
		if strings != nil {
			self.NodeTemplateName = &(*strings)[0]
			self.InterfaceName = &(*strings)[1]
		}
	}
	return self
}

// tosca.Renderable interface
func (self *InterfaceMapping) Render() {
	log.Info("{render} interface mapping")

	if (self.NodeTemplate == nil) || (self.InterfaceName == nil) {
		return
	}

	name := *self.InterfaceName
	if _, ok := self.NodeTemplate.Interfaces[name]; !ok {
		self.Context.ListChild(1, name).ReportReferenceNotFound("interface", self.NodeTemplate)
	}
}
