package v1_1

import (
	"github.com/tliron/puccini/tosca"
)

//
// RelationshipDefinition
//

type RelationshipDefinition struct {
	*Entity `name:"relationship definition"`

	RelationshipTypeName *string              `read:"type"` // required only if cannot be inherited
	InterfaceDefinitions InterfaceDefinitions `read:"interfaces,InterfaceDefinition" inherit:"interfaces,RelationshipType"`

	RelationshipType *RelationshipType `lookup:"type,RelationshipTypeName" json:"-" yaml:"-"`

	typeMissingProblemReported bool
}

func NewRelationshipDefinition(context *tosca.Context) *RelationshipDefinition {
	return &RelationshipDefinition{
		Entity:               NewEntity(context),
		InterfaceDefinitions: make(InterfaceDefinitions),
	}
}

// tosca.Reader signature
func ReadRelationshipDefinition(context *tosca.Context) interface{} {
	self := NewRelationshipDefinition(context)
	if context.Is("map") {
		context.ValidateUnsupportedFields(context.ReadFields(self, Readers))
	} else if context.ValidateType("map", "string") {
		self.RelationshipTypeName = context.ReadString()
	}
	return self
}

func (self *RelationshipDefinition) NewDefaultAssignment(context *tosca.Context) *RelationshipAssignment {
	assignment := NewRelationshipAssignment(context)
	assignment.RelationshipTemplateOrRelationshipTypeName = self.RelationshipTypeName
	assignment.RelationshipType = self.RelationshipType
	return assignment
}

func (self *RelationshipDefinition) Inherit(parentDefinition *RelationshipDefinition) {
	if parentDefinition != nil {
		if (self.RelationshipTypeName == nil) && (parentDefinition.RelationshipTypeName != nil) {
			self.RelationshipTypeName = parentDefinition.RelationshipTypeName
		}
		if (self.RelationshipType == nil) && (parentDefinition.RelationshipType != nil) {
			self.RelationshipType = parentDefinition.RelationshipType
		}

		// Validate type compatibility
		if (self.RelationshipType != nil) && (parentDefinition.RelationshipType != nil) && !self.Context.Hierarchy.IsCompatible(parentDefinition.RelationshipType, self.RelationshipType) {
			self.Context.ReportIncompatibleType(self.RelationshipType.Name, parentDefinition.RelationshipType.Name)
			return
		}

		self.InterfaceDefinitions.Inherit(parentDefinition.InterfaceDefinitions)
	} else {
		self.InterfaceDefinitions.Inherit(nil)
	}

	if self.RelationshipTypeName == nil {
		// Avoid reporting more than once
		if !self.typeMissingProblemReported {
			self.Context.FieldChild("type", nil).ReportFieldMissing()
			self.typeMissingProblemReported = true
		}
	}
}
