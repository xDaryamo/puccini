package tosca_v2_0

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
)

//
// NotificationDefinition
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.19
//

type NotificationDefinition struct {
	*Entity `name:"notification definition"`
	Name    string

	Description    *string                  `read:"description"`
	Implementation *InterfaceImplementation `read:"implementation,InterfaceImplementation"`
	Outputs        OutputMappings           `read:"outputs,OutputMapping"`
}

func NewNotificationDefinition(context *tosca.Context) *NotificationDefinition {
	return &NotificationDefinition{
		Entity:  NewEntity(context),
		Name:    context.Name,
		Outputs: make(OutputMappings),
	}
}

// tosca.Reader signature
func ReadNotificationDefinition(context *tosca.Context) tosca.EntityPtr {
	self := NewNotificationDefinition(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
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
	logInherit.Debugf("notification definition: %s", self.Name)

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
