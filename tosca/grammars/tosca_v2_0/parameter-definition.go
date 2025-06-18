package tosca_v2_0

import (
	"errors"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// ParameterDefinition
//
// [TOSCA-v2.0] @ 9.8
//

type ParameterDefinition struct {
	*PropertyDefinition `name:"parameter definition"`

	Value         *Value  `read:"value,Value"`
	NodeTemplate  *string `traverse:"ignore" json:"-" yaml:"-"` // parsed from mapping
	AttributeName *string `traverse:"ignore" json:"-" yaml:"-"` // parsed from mapping
}

func NewParameterDefinition(context *parsing.Context) *ParameterDefinition {
	return &ParameterDefinition{PropertyDefinition: NewPropertyDefinition(context)}
}

// ([parsing.Reader] signature)
func ReadParameterDefinition(context *parsing.Context) parsing.EntityPtr {
	self := NewParameterDefinition(context)
	self.looseType = true // type is optional for parameters

	// Check for short-hand mapping syntax: [ SELF, attribute_name ]
	if context.Is(ard.TypeList) {
		if list := context.Data.(ard.List); len(list) >= 2 {
			if nodeTemplate, ok := list[0].(string); ok {
				if attributeName, ok := list[1].(string); ok {
					self.NodeTemplate = &nodeTemplate
					self.AttributeName = &attributeName
					return self
				}
			}
		}
		context.ReportValueMalformed("parameter definition", "mapping list must have at least 2 string elements")
		return self
	}

	// Standard map notation
	if context.Is(ard.TypeMap) {
		// Manually read the mapping field if present BEFORE reading other fields
		if data, ok := context.Data.(ard.Map); ok {
			if mappingData, exists := data["mapping"]; exists {
				// Parse mapping as a list of strings
				if mappingList, ok := mappingData.(ard.List); ok && len(mappingList) >= 2 {
					if nodeTemplate, ok := mappingList[0].(string); ok {
						if attributeName, ok := mappingList[1].(string); ok {
							self.NodeTemplate = &nodeTemplate
							self.AttributeName = &attributeName
						}
					}
				} else {
					context.FieldChild("mapping", mappingData).ReportValueMalformed("mapping", "must be a list of at least 2 strings")
				}
			}
		}

		// Read the standard PropertyDefinition fields
		var ignore []string
		if context.HasQuirk(parsing.QuirkAnnotationsIgnore) {
			ignore = append(ignore, "annotations")
		}
		// Add mapping to ignored fields since we handle it manually
		ignore = append(ignore, "mapping")

		context.ValidateUnsupportedFields(append(context.ReadFields(self), ignore...))

		// Validate mutual exclusivity
		if (self.Value != nil) && (self.NodeTemplate != nil && self.AttributeName != nil) {
			context.ReportError(errors.New("'value' and 'mapping' are mutually exclusive"))
		}
	} else if context.ValidateType(ard.TypeString, ard.TypeInteger, ard.TypeFloat, ard.TypeBoolean) {
		// Single value (for outgoing parameters)
		self.Value = ReadValue(context).(*Value)
	} else {
		context.ReportValueWrongType(ard.TypeMap, ard.TypeList)
	}

	return self
}

func (self *ParameterDefinition) Inherit(parentDefinition *ParameterDefinition) {
	logInherit.Debugf("parameter definition: %s", self.Name)

	self.PropertyDefinition.Inherit(parentDefinition.PropertyDefinition)

	if (self.Value == nil) && (parentDefinition.Value != nil) {
		self.Value = parentDefinition.Value
	}

	if (self.NodeTemplate == nil) && (parentDefinition.NodeTemplate != nil) {
		self.NodeTemplate = parentDefinition.NodeTemplate
	}

	if (self.AttributeName == nil) && (parentDefinition.AttributeName != nil) {
		self.AttributeName = parentDefinition.AttributeName
	}
}

// Check if this is an outgoing parameter (has value)
func (self *ParameterDefinition) IsOutgoing() bool {
	return self.Value != nil
}

// Check if this is an incoming parameter (has mapping)
func (self *ParameterDefinition) IsIncoming() bool {
	return (self.NodeTemplate != nil) && (self.AttributeName != nil)
}

// ([parsing.Renderable] interface)
func (self *ParameterDefinition) Render() {
	// Avoid rendering more than once
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

func (self *ParameterDefinition) Normalize(context *parsing.Context) normal.Value {
	var value *Value
	if self.Value != nil {
		value = self.Value
	} else {
		// Parameters should always appear, even if they have no default value
		value = NewValue(context.MapChild(self.Name, nil))
	}
	return value.Normalize()
}

// Apply parameter mapping to node attributes AND properties
func (self *ParameterDefinition) ApplyMapping(nodeTemplates NodeTemplates, inputValue ard.Value) {
	if self.IsIncoming() {
		nodeName := *self.NodeTemplate
		attributeName := *self.AttributeName

		// Find the node template
		for _, nodeTemplate := range nodeTemplates {
			if nodeTemplate.Name == nodeName {
				// Create a Value from the input
				valueContext := nodeTemplate.Context.FieldChild("properties", nil).FieldChild(attributeName, inputValue)
				value := NewValue(valueContext)
				value.Context.Data = inputValue

				// Set both as property AND attribute
				// Properties (for validation and requirements)
				if nodeTemplate.Properties == nil {
					nodeTemplate.Properties = make(Values)
				}
				nodeTemplate.Properties[attributeName] = value

				// Attributes (for runtime access)
				if nodeTemplate.Attributes == nil {
					nodeTemplate.Attributes = make(Values)
				}
				nodeTemplate.Attributes[attributeName] = value
				break
			}
		}
	}
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

func (self ParameterDefinitions) Normalize(c normal.Values, context *parsing.Context) {
	for key, definition := range self {
		c[key] = definition.Normalize(context)
	}
}

// Apply all parameter mappings to node templates
func (self ParameterDefinitions) ApplyMappings(nodeTemplates NodeTemplates, inputs map[string]ard.Value) {
	for name, definition := range self {
		if definition.IsIncoming() {
			if inputValue, exists := inputs[name]; exists {
				definition.ApplyMapping(nodeTemplates, inputValue)
			}
		}
	}
}
