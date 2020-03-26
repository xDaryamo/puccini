package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// NotificationDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.19
//

type NotificationDefinition struct {
	*Entity `name:"notification definition"`
	Name    string

	Description    *string                  `read:"description"`
	Implementation *InterfaceImplementation `read:"implementation,InterfaceImplementation"`
	Outputs        NotificationOutputs      `read:"outputs,NotificationOutput"`
}

func NewNotificationDefinition(context *tosca.Context) *NotificationDefinition {
	return &NotificationDefinition{
		Entity:  NewEntity(context),
		Name:    context.Name,
		Outputs: make(NotificationOutputs),
	}
}

// tosca.Reader signature
func ReadNotificationDefinition(context *tosca.Context) interface{} {
	self := NewNotificationDefinition(context)

	if context.Is("!!map") {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType("!!map", "!!str") {
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
	log.Infof("{inherit} notification definition: %s", self.Name)

	if (self.Description == nil) && (parentDefinition.Description != nil) {
		self.Description = parentDefinition.Description
	}

	self.Outputs.Inherit(parentDefinition.Outputs)
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
		}
	}
}
