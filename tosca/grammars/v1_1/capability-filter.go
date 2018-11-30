package v1_1

import (
	"github.com/tliron/puccini/tosca"
)

//
// CapabilityFilter
//
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.6
//

type CapabilityFilter struct {
	*Entity `name:"capability filter"`
	Name    string

	Properties PropertyFilters `read:"properties,PropertyFilter"`
}

func NewCapabilityFilter(context *tosca.Context) *CapabilityFilter {
	return &CapabilityFilter{
		Entity:     NewEntity(context),
		Name:       context.Name,
		Properties: make(PropertyFilters),
	}
}

// tosca.Reader signature
func ReadCapabilityFilter(context *tosca.Context) interface{} {
	self := NewCapabilityFilter(context)
	context.ValidateUnsupportedFields(context.ReadFields(self, Readers))
	return self
}

// tosca.Mappable interface
func (self *CapabilityFilter) GetKey() string {
	return self.Name
}

//
// CapabilityFilters
//

type CapabilityFilters map[string]*CapabilityFilter
