package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
)

//
// PropertyDefinition
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.10
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.9
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.8
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.8
//

type PropertyDefinition struct {
	*AttributeDefinition `name:"property definition"`

	Required          *bool             `read:"required"`
	ConstraintClauses ConstraintClauses `read:"constraints,[]ConstraintClause"`
}

func NewPropertyDefinition(context *tosca.Context) *PropertyDefinition {
	return &PropertyDefinition{AttributeDefinition: NewAttributeDefinition(context)}
}

// tosca.Reader signature
func ReadPropertyDefinition(context *tosca.Context) tosca.EntityPtr {
	self := NewPropertyDefinition(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

func (self *PropertyDefinition) Inherit(parentDefinition *PropertyDefinition) {
	logInherit.Debugf("property definition: %s", self.Name)

	self.AttributeDefinition.Inherit(parentDefinition.AttributeDefinition)

	if (self.Required == nil) && (parentDefinition.Required != nil) {
		self.Required = parentDefinition.Required
	}
	if parentDefinition.ConstraintClauses != nil {
		self.ConstraintClauses = parentDefinition.ConstraintClauses.Append(self.ConstraintClauses)
	}
}

// parser.Renderable interface
func (self *PropertyDefinition) Render() {
	logRender.Debugf("property definition: %s", self.Name)

	self.AttributeDefinition.Render()
	self.ConstraintClauses.Render(self.DataType)
}

func (self *PropertyDefinition) IsRequired() bool {
	// defaults to true
	return (self.Required == nil) || *self.Required
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
		if parentDefinition, ok := parentDefinitions[name]; ok {
			if definition != parentDefinition {
				definition.Inherit(parentDefinition)
			}
		}
	}
}
