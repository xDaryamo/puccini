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

func (self *PropertyDefinition) Inherit(parentDefinition *PropertyDefinition) {
	if parentDefinition != nil {
		self.ParameterDefinition.Inherit(parentDefinition.ParameterDefinition)

		if (self.Required == nil) && (parentDefinition.Required != nil) {
			self.Required = parentDefinition.Required
		}
	} else {
		self.ParameterDefinition.Inherit(nil)
	}
}

//
// PropertyDefinitions
//

type PropertyDefinitions map[string]*PropertyDefinition

func (self PropertyDefinitions) Inherit(parentDefinitions PropertyDefinitions) {
	for name, definition := range parentDefinitions {
		if _, ok := self[name]; !ok {
			self[name] = definition
		}
	}

	for name, definition := range self {
		if parentDefinitions != nil {
			if parentDefinition, ok := parentDefinitions[name]; ok {
				if definition != parentDefinition {
					definition.Inherit(parentDefinition)
				}
				continue
			}
		}

		definition.Inherit(nil)
	}
}
