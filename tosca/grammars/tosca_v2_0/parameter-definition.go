package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// ParameterDefinition
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.14
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.13
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.12
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.12
//

type ParameterDefinition struct {
	*PropertyDefinition `name:"parameter definition"`

	Value *Value `read:"value,Value"`
}

func NewParameterDefinition(context *tosca.Context) *ParameterDefinition {
	return &ParameterDefinition{PropertyDefinition: NewPropertyDefinition(context)}
}

// tosca.Reader signature
func ReadParameterDefinition(context *tosca.Context) tosca.EntityPtr {
	self := NewParameterDefinition(context)
	self.looseType = true
	var ignore []string
	if context.HasQuirk(tosca.QuirkAnnotationsIgnore) {
		ignore = append(ignore, "annotations")
	}
	context.ValidateUnsupportedFields(append(context.ReadFields(self), ignore...))
	return self
}

func (self *ParameterDefinition) Inherit(parentDefinition *ParameterDefinition) {
	logInherit.Debugf("parameter definition: %s", self.Name)

	self.PropertyDefinition.Inherit(parentDefinition.PropertyDefinition)

	if (self.Value == nil) && (parentDefinition.Value != nil) {
		self.Value = parentDefinition.Value
	}
}

// tosca.Renderable interface
func (self *ParameterDefinition) Render() {
	self.renderOnce.Do(self.render)
}

func (self *ParameterDefinition) render() {
	logRender.Debugf("parameter definition: %s", self.Name)

	self.PropertyDefinition.render()

	if self.Value != nil {
		if self.DataType != nil {
			self.Value.RenderProperty(self.DataType, self.PropertyDefinition)
		}
	} else if self.Default != nil {
		self.Value = self.Default
	}
}

func (self *ParameterDefinition) Normalize(context *tosca.Context) normal.Value {
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

func (self ParameterDefinitions) Inherit(parentDefinitions ParameterDefinitions) {
	for name, definition := range parentDefinitions {
		if _, ok := self[name]; !ok {
			self[name] = definition
		}
	}

	for name, definition := range self {
		if parentDefinition, ok := parentDefinitions[name]; ok {
			if definition != parentDefinition {
				definition.Inherit(parentDefinition)
			}
		}
	}
}

func (self ParameterDefinitions) Render(kind string, mapped []string) {
	for _, definition := range self {
		definition.Render()

		if definition.Value == nil {
			isMapped := false
			for _, mapped_ := range mapped {
				if definition.Name == mapped_ {
					isMapped = true
					break
				}
			}

			if !isMapped && definition.IsRequired() {
				definition.Context.ReportValueRequired(kind)
				return
			}
		}
	}
}

func (self ParameterDefinitions) Normalize(c normal.Values, context *tosca.Context) {
	for key, definition := range self {
		c[key] = definition.Normalize(context)
	}
}
