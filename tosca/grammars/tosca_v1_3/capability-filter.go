package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// CapabilityFilter
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.5.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.4.2
//

type CapabilityFilter struct {
	*Entity `name:"capability filter"`
	Name    string

	PropertyFilters PropertyFilters `read:"properties,PropertyFilter"`
}

func NewCapabilityFilter(context *tosca.Context) *CapabilityFilter {
	return &CapabilityFilter{
		Entity:          NewEntity(context),
		Name:            context.Name,
		PropertyFilters: make(PropertyFilters),
	}
}

// tosca.Reader signature
func ReadCapabilityFilter(context *tosca.Context) interface{} {
	self := NewCapabilityFilter(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Mappable interface
func (self *CapabilityFilter) GetKey() string {
	return self.Name
}

func (self CapabilityFilter) Normalize(r *normal.Requirement) normal.FunctionCallMap {
	if len(self.PropertyFilters) == 0 {
		return nil
	}

	functionCallMap := make(normal.FunctionCallMap)
	r.CapabilityPropertyConstraints[self.Name] = functionCallMap
	self.PropertyFilters.Normalize(functionCallMap)

	return functionCallMap
}

//
// CapabilityFilters
//

type CapabilityFilters map[string]*CapabilityFilter

func (self CapabilityFilters) Normalize(r *normal.Requirement) {
	for _, capabilityFilter := range self {
		capabilityFilter.Normalize(r)
	}
}
