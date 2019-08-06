package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// ParameterDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.13
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.12
//

type ParameterDefinition struct {
	*PropertyDefinition `name:"parameter definition"`

	Value *Value `read:"value,Value"`
}

func NewParameterDefinition(context *tosca.Context) *ParameterDefinition {
	return &ParameterDefinition{PropertyDefinition: NewPropertyDefinition(context)}
}

// tosca.Reader signature
func ReadParameterDefinition(context *tosca.Context) interface{} {
	self := NewParameterDefinition(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

func (self *ParameterDefinition) Render(kind string) {
	// TODO: what to do if there is no "type"?

	if self.Value == nil {
		self.Value = self.Default
	}

	if self.Value == nil {
		// PropertyDefinition.Required defaults to true
		required := (self.Required == nil) || *self.Required
		if required {
			self.Context.ReportPropertyRequired(kind)
			return
		}
	} else if self.DataType != nil {
		self.Value.RenderProperty(self.DataType, self.PropertyDefinition)
	}
}

func (self *ParameterDefinition) Normalize(context *tosca.Context) normal.Constrainable {
	var value *Value
	if self.Value != nil {
		value = self.Value
	} else {
		// Parameters should always appear, even if they have no default value
		value = NewValue(context.MapChild(self.Name, nil))
	}
	return value.Normalize()
}

//
// ParameterDefinitions
//

type ParameterDefinitions map[string]*ParameterDefinition

func (self ParameterDefinitions) Render(kind string, context *tosca.Context) {
	for _, definition := range self {
		definition.Render(kind)
	}
}

func (self ParameterDefinitions) Normalize(c normal.Constrainables, context *tosca.Context) {
	for key, definition := range self {
		c[key] = definition.Normalize(context)
	}
}
