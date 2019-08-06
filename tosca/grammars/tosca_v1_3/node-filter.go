package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// NodeFilter
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.5
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.4
//

type NodeFilter struct {
	*Entity `name:"node filter"`

	PropertyFilters   PropertyFilters   `read:"properties,PropertyFilter"`
	CapabilityFilters CapabilityFilters `read:"capabilities,CapabilityFilter"`
}

func NewNodeFilter(context *tosca.Context) *NodeFilter {
	return &NodeFilter{
		Entity:            NewEntity(context),
		PropertyFilters:   make(PropertyFilters),
		CapabilityFilters: make(CapabilityFilters),
	}
}

// tosca.Reader signature
func ReadNodeFilter(context *tosca.Context) interface{} {
	self := NewNodeFilter(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

func (self *NodeFilter) Normalize(r *normal.Requirement) {
	self.PropertyFilters.Normalize(r.NodeTemplatePropertyConstraints)
	self.CapabilityFilters.Normalize(r)
}
