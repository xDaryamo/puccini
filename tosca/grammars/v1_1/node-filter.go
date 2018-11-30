package v1_1

import (
	"github.com/tliron/puccini/tosca"
)

//
// NodeFilter
//
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.4
//

type NodeFilter struct {
	*Entity `name:"node filter"`

	Properties        PropertyFilters   `read:"properties,PropertyFilter"`
	CapabilityFilters CapabilityFilters `read:"capabilities,CapabilityFilter"`
}

func NewNodeFilter(context *tosca.Context) *NodeFilter {
	return &NodeFilter{
		Entity:            NewEntity(context),
		Properties:        make(PropertyFilters),
		CapabilityFilters: make(CapabilityFilters),
	}
}

// tosca.Reader signature
func ReadNodeFilter(context *tosca.Context) interface{} {
	self := NewNodeFilter(context)
	context.ValidateUnsupportedFields(context.ReadFields(self, Readers))
	return self
}

func (self *NodeFilter) FilterNodeTemplates(nodeTemplates []*NodeTemplate) []*NodeTemplate {
	if len(self.Properties) == 0 {
		return nodeTemplates
	}

	var filteredNodeTemplates []*NodeTemplate
	for _, nodeTemplate := range nodeTemplates {
		if self.Properties.Apply(nodeTemplate) {
			filteredNodeTemplates = append(filteredNodeTemplates, nodeTemplate)
		}
	}

	return filteredNodeTemplates
}
