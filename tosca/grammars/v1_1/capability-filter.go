package v1_1

import (
	"github.com/tliron/puccini/tosca"
)

//
// CapabilityFilter
//

type CapabilityFilter struct {
	*Entity `name:"capability filter"`
	Name    string

	Properties Values `read:"properties,Value"`
}

func NewCapabilityFilter(context *tosca.Context) *CapabilityFilter {
	return &CapabilityFilter{
		Entity:     NewEntity(context),
		Name:       context.Name,
		Properties: make(Values),
	}
}

// tosca.Reader signature
func ReadCapabilityFilter(context *tosca.Context) interface{} {
	self := NewCapabilityFilter(context)
	context.ValidateUnsupportedFields(context.ReadFields(self, Readers))
	return self
}

func init() {
	Readers["CapabilityFilter"] = ReadCapabilityFilter
}

// tosca.Mappable interface
func (self *CapabilityFilter) GetKey() string {
	return self.Name
}

//
// CapabilityFilters
//

type CapabilityFilters map[string]*CapabilityFilter
