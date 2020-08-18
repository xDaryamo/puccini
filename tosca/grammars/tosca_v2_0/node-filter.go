package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// NodeFilter
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.5
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.5
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.4
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.4
//

type NodeFilter struct {
	*Entity `name:"node filter"`

	PropertyFilters   PropertyFilters   `read:"properties,PropertyFilter"`
	CapabilityFilters CapabilityFilters `read:"capabilities,{}CapabilityFilter"`
}

func NewNodeFilter(context *tosca.Context) *NodeFilter {
	return &NodeFilter{
		Entity:          NewEntity(context),
		PropertyFilters: make(PropertyFilters),
	}
}

// tosca.Reader signature
func ReadNodeFilter(context *tosca.Context) tosca.EntityPtr {
	self := NewNodeFilter(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

func (self *NodeFilter) Normalize(normalRequirement *normal.Requirement) {
	self.PropertyFilters.Normalize(normalRequirement.NodeTemplatePropertyConstraints)
	self.CapabilityFilters.Normalize(normalRequirement)
}
