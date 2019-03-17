package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// PropertyDefinition
//

type PropertyDefinition struct {
	*ParameterDefinition `name:"property definition"`

	Required *bool `read:"required"`
}

func NewPropertyDefinition(context *tosca.Context) *PropertyDefinition {
	return &PropertyDefinition{
		ParameterDefinition: NewParameterDefinition(context),
	}
}

// tosca.Reader signature
func ReadPropertyDefinition(context *tosca.Context) interface{} {
	self := NewPropertyDefinition(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Mappable interface
func (self *PropertyDefinition) GetKey() string {
	return self.Name
}

//
// PropertyDefinitions
//

type PropertyDefinitions map[string]*PropertyDefinition
