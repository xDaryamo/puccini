package tosca_v2_0

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
)

//
// InterfaceMapping
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.12
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.11
//

type InterfaceMapping struct {
	*Entity `name:"interface mapping"`

	NodeTemplateName *string `require:"0"`
	InterfaceName    *string `require:"1"`

	NodeTemplate *NodeTemplate `lookup:"0,NodeTemplateName" json:"-" yaml:"-"`
}

func NewInterfaceMapping(context *tosca.Context) *InterfaceMapping {
	return &InterfaceMapping{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadInterfaceMapping(context *tosca.Context) tosca.EntityPtr {
	self := NewInterfaceMapping(context)
	if context.ValidateType(ard.TypeList) {
		strings := context.ReadStringListFixed(2)
		if strings != nil {
			self.NodeTemplateName = &(*strings)[0]
			self.InterfaceName = &(*strings)[1]
		}
	}
	return self
}

// parser.Renderable interface
func (self *InterfaceMapping) Render() {
	logRender.Debug("interface mapping")

	if (self.NodeTemplate == nil) || (self.InterfaceName == nil) {
		return
	}

	name := *self.InterfaceName
	self.NodeTemplate.Render()
	if _, ok := self.NodeTemplate.Interfaces[name]; !ok {
		self.Context.ListChild(1, name).ReportReferenceNotFound("interface", self.NodeTemplate)
	}
}

//
// InterfaceMappings
//

type InterfaceMappings []*InterfaceMapping
