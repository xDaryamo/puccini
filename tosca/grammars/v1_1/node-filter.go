package v1_1

import (
	"github.com/tliron/puccini/tosca"
)

//
// NodeFilter
//

type NodeFilter struct {
	*Entity `name:"node filter"`

	Properties        Values            `read:"properties,Value"`
	CapabilityFilters CapabilityFilters `read:"capabilities,CapabilityFilter"`
}

func NewNodeFilter(context *tosca.Context) *NodeFilter {
	return &NodeFilter{
		Entity:            NewEntity(context),
		Properties:        make(Values),
		CapabilityFilters: make(CapabilityFilters),
	}
}

// tosca.Reader signature
func ReadNodeFilter(context *tosca.Context) interface{} {
	self := NewNodeFilter(context)
	context.ValidateUnsupportedFields(context.ReadFields(self, Readers))
	return self
}
