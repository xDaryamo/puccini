package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// OperationDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.15
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.13
//

type OperationDefinition struct {
	*Entity `name:"operation definition"`
	Name    string

	Description      *string                  `read:"description"`
	Implementation   *OperationImplementation `read:"implementation,OperationImplementation"`
	InputDefinitions PropertyDefinitions      `read:"inputs,PropertyDefinition"`
}

func NewOperationDefinition(context *tosca.Context) *OperationDefinition {
	return &OperationDefinition{
		Entity:           NewEntity(context),
		Name:             context.Name,
		InputDefinitions: make(PropertyDefinitions),
	}
}

// tosca.Reader signature
func ReadOperationDefinition(context *tosca.Context) interface{} {
	self := NewOperationDefinition(context)

	if context.Is("map") {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType("map", "string") {
		// Short notation
		self.Implementation = ReadOperationImplementation(context.FieldChild("implementation", context.Data)).(*OperationImplementation)
	}

	return self
}

// tosca.Mappable interface
func (self *OperationDefinition) GetKey() string {
	return self.Name
}

func (self *OperationDefinition) Inherit(parentDefinition *OperationDefinition) {
	if parentDefinition != nil {
		if (self.Description == nil) && (parentDefinition.Description != nil) {
			self.Description = parentDefinition.Description
		}

		self.InputDefinitions.Inherit(parentDefinition.InputDefinitions)
	} else {
		self.InputDefinitions.Inherit(nil)
	}
}

func (self *OperationDefinition) Normalize(o *normal.Operation) {
	if self.Description != nil {
		o.Description = *self.Description
	}

	if self.Implementation != nil {
		self.Implementation.Normalize(o)
	}

	// TODO: input definitions
	//self.InputsDefinitions.Normalize(o.Inputs)
}

//
// OperationDefinitions
//

type OperationDefinitions map[string]*OperationDefinition

func (self OperationDefinitions) Inherit(parentDefinitions OperationDefinitions) {
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
		} else {
			definition.Inherit(nil)
		}
	}
}
