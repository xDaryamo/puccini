package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// NotificationAssignment
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.19
//

type NotificationAssignment struct {
	*Entity `name:"notification"`
	Name    string

	Description    *string                  `read:"description"`
	Implementation *InterfaceImplementation `read:"implementation,InterfaceImplementation"`
	Outputs        AttributeMappings        `read:"outputs,AttributeMapping"`
}

func NewNotificationAssignment(context *tosca.Context) *NotificationAssignment {
	return &NotificationAssignment{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadNotificationAssignment(context *tosca.Context) interface{} {
	self := NewNotificationAssignment(context)

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
func (self *NotificationAssignment) GetKey() string {
	return self.Name
}

//
// NotificationAssignments
//

type NotificationAssignments map[string]*NotificationAssignment
