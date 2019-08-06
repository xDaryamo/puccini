package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// InterfaceDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.16
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.14
//

type InterfaceDefinition struct {
	*Entity `name:"interface definition" json:"-" yaml:"-"`
	Name    string

	InterfaceTypeName    *string              `read:"type"` // required only if cannot be inherited
	InputDefinitions     PropertyDefinitions  `read:"inputs,PropertyDefinition" inherit:"inputs,InterfaceType"`
	OperationDefinitions OperationDefinitions `read:"?,OperationDefinition" inherit:"?,InterfaceType"`

	InterfaceType *InterfaceType `lookup:"type,InterfaceTypeName" json:"-" yaml:"-"`

	typeMissingProblemReported bool
}

func NewInterfaceDefinition(context *tosca.Context) *InterfaceDefinition {
	return &InterfaceDefinition{
		Entity:               NewEntity(context),
		Name:                 context.Name,
		InputDefinitions:     make(PropertyDefinitions),
		OperationDefinitions: make(OperationDefinitions),
	}
}

// tosca.Reader signature
func ReadInterfaceDefinition(context *tosca.Context) interface{} {
	self := NewInterfaceDefinition(context)
	context.ReadFields(self)
	return self
}

// tosca.Mappable interface
func (self *InterfaceDefinition) GetKey() string {
	return self.Name
}

func (self *InterfaceDefinition) Inherit(parentDefinition *InterfaceDefinition) {
	if parentDefinition != nil {
		if (self.InterfaceTypeName == nil) && (parentDefinition.InterfaceTypeName != nil) {
			self.InterfaceTypeName = parentDefinition.InterfaceTypeName
		}
		if (self.InterfaceType == nil) && (parentDefinition.InterfaceType != nil) {
			self.InterfaceType = parentDefinition.InterfaceType
		}

		// Validate type compatibility
		if (self.InterfaceType != nil) && (parentDefinition.InterfaceType != nil) && !self.Context.Hierarchy.IsCompatible(parentDefinition.InterfaceType, self.InterfaceType) {
			self.Context.ReportIncompatibleType(self.InterfaceType.Name, parentDefinition.InterfaceType.Name)
			return
		}

		self.InputDefinitions.Inherit(parentDefinition.InputDefinitions)
		self.OperationDefinitions.Inherit(parentDefinition.OperationDefinitions)
	} else {
		self.InputDefinitions.Inherit(nil)
		self.OperationDefinitions.Inherit(nil)
	}

	if self.InterfaceTypeName == nil {
		// Avoid reporting more than once
		if !self.typeMissingProblemReported {
			self.Context.FieldChild("type", nil).ReportFieldMissing()
			self.typeMissingProblemReported = true
		}
	}
}

//
// InterfaceDefinitions
//

type InterfaceDefinitions map[string]*InterfaceDefinition

func (self InterfaceDefinitions) Inherit(parentDefinitions InterfaceDefinitions) {
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
