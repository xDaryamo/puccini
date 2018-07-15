package v1_1

import (
	"github.com/tliron/puccini/tosca"
)

//
// RequirementDefinition
//

type RequirementDefinition struct {
	*Entity `name:"requirement definition"`
	Name    string

	NodeTypeName           *string                 `read:"node"`
	CapabilityTypeName     *string                 `read:"capability"` // required only if cannot be inherited
	RelationshipDefinition *RelationshipDefinition `read:"relationship,RelationshipDefinition"`
	Occurrences            *RangeEntity            `read:"occurrences,RangeEntity"`

	NodeType       *NodeType       `lookup:"node,NodeTypeName" json:"-" yaml:"-"`
	CapabilityType *CapabilityType `lookup:"capability,CapabilityTypeName" json:"-" yaml:"-"`

	capabilityMissingProblemReported bool
}

func NewRequirementDefinition(context *tosca.Context) *RequirementDefinition {
	return &RequirementDefinition{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadRequirementDefinition(context *tosca.Context) interface{} {
	self := NewRequirementDefinition(context)
	if context.Is("map") {
		context.ValidateUnsupportedFields(context.ReadFields(self, Readers))
	} else if context.ValidateType("map", "string") {
		self.CapabilityTypeName = context.ReadString()
	}
	return self
}

// tosca.Mappable interface
func (self *RequirementDefinition) GetKey() string {
	return self.Name
}

func (self *RequirementDefinition) Inherit(parentDefinition *RequirementDefinition) {
	if parentDefinition != nil {
		if (self.NodeTypeName == nil) && (parentDefinition.NodeTypeName != nil) {
			self.NodeTypeName = parentDefinition.NodeTypeName
		}
		if (self.CapabilityTypeName == nil) && (parentDefinition.CapabilityTypeName != nil) {
			self.CapabilityTypeName = parentDefinition.CapabilityTypeName
		}
		if (self.RelationshipDefinition == nil) && (parentDefinition.RelationshipDefinition != nil) {
			self.RelationshipDefinition = parentDefinition.RelationshipDefinition
		}
		if (self.Occurrences == nil) && (parentDefinition.Occurrences != nil) {
			self.Occurrences = parentDefinition.Occurrences
		}
		if (self.NodeType == nil) && (parentDefinition.NodeType != nil) {
			self.NodeType = parentDefinition.NodeType
		}
		if (self.CapabilityType == nil) && (parentDefinition.CapabilityType != nil) {
			self.CapabilityType = parentDefinition.CapabilityType
		}

		// Validate type compatibility
		if (self.NodeType != nil) && (parentDefinition.NodeType != nil) && !self.Context.Hierarchy.IsCompatible(parentDefinition.NodeType, self.NodeType) {
			self.Context.ReportIncompatibleType(self.NodeType.Name, parentDefinition.NodeType.Name)
		}
		if (self.CapabilityType != nil) && (parentDefinition.CapabilityType != nil) && !self.Context.Hierarchy.IsCompatible(parentDefinition.CapabilityType, self.CapabilityType) {
			self.Context.ReportIncompatibleType(self.CapabilityType.Name, parentDefinition.CapabilityType.Name)
		}
	}

	if self.CapabilityTypeName == nil {
		// Avoid reporting more than once
		if !self.capabilityMissingProblemReported {
			self.Context.FieldChild("capability", nil).ReportFieldMissing()
			self.capabilityMissingProblemReported = true
		}
	}

	if self.RelationshipDefinition != nil {
		if parentDefinition != nil {
			self.RelationshipDefinition.Inherit(parentDefinition.RelationshipDefinition)
		} else {
			self.RelationshipDefinition.Inherit(nil)
		}
	}
}

//
// RequirementDefinitions
//

type RequirementDefinitions map[string]*RequirementDefinition

func (self RequirementDefinitions) Inherit(parentDefinitions RequirementDefinitions) {
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
