package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// NotificationDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.19
//

type NotificationDefinition struct {
	*Entity `name:"operation definition"`
	Name    string

	Description       *string                  `read:"description"`
	Implementation    *InterfaceImplementation `read:"implementation,InterfaceImplementation"`
	OutputDefinitions AttributeMappings        `read:"outputs,AttributeMapping"`
}

func NewNotificationDefinition(context *tosca.Context) *NotificationDefinition {
	return &NotificationDefinition{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadNotificationDefinition(context *tosca.Context) interface{} {
	self := NewNotificationDefinition(context)

	if context.Is("map") {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType("map", "string") {
		// Short notation
		self.Implementation = ReadInterfaceImplementation(context.FieldChild("implementation", context.Data)).(*InterfaceImplementation)
	}

	return self
}

// tosca.Mappable interface
func (self *NotificationDefinition) GetKey() string {
	return self.Name
}

func (self *NotificationDefinition) Inherit(parentDefinition *NotificationDefinition) {
	if parentDefinition != nil {
		if (self.Description == nil) && (parentDefinition.Description != nil) {
			self.Description = parentDefinition.Description
		}

		// TODO
		//self.OutputDefinitions.Inherit(parentDefinition.OutputDefinitions)
	} else {
		//self.OutputDefinitions.Inherit(nil)
	}
}

func (self *NotificationDefinition) Normalize(o *normal.Operation) {
	if self.Description != nil {
		o.Description = *self.Description
	}

	if self.Implementation != nil {
		self.Implementation.Normalize(o)
	}

	// TODO: output definitions
	//self.OutputDefinitions.Normalize(o.Inputs)
}

//
// NotificationDefinitions
//

type NotificationDefinitions map[string]*NotificationDefinition

func (self NotificationDefinitions) Inherit(parentDefinitions NotificationDefinitions) {
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
