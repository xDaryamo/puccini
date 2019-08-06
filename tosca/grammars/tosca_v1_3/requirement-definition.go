package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// RequirementDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.3
//

type RequirementDefinition struct {
	*Entity `name:"requirement definition"`
	Name    string

	TargetCapabilityTypeName *string                 `read:"capability"` // required only if cannot be inherited
	TargetNodeTypeName       *string                 `read:"node"`
	RelationshipDefinition   *RelationshipDefinition `read:"relationship,RelationshipDefinition"`
	Occurrences              *RangeEntity            `read:"occurrences,RangeEntity"`

	TargetCapabilityType *CapabilityType `lookup:"capability,TargetCapabilityTypeName" json:"-" yaml:"-"`
	TargetNodeType       *NodeType       `lookup:"node,TargetNodeTypeName" json:"-" yaml:"-"`

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
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType("map", "string") {
		// Short notation
		self.TargetCapabilityTypeName = context.FieldChild("capability", context.Data).ReadString()
	}

	return self
}

// tosca.Mappable interface
func (self *RequirementDefinition) GetKey() string {
	return self.Name
}

func (self *RequirementDefinition) Inherit(parentDefinition *RequirementDefinition) {
	if parentDefinition != nil {
		if (self.TargetCapabilityTypeName == nil) && (parentDefinition.TargetCapabilityTypeName != nil) {
			self.TargetCapabilityTypeName = parentDefinition.TargetCapabilityTypeName
		}
		if (self.TargetNodeTypeName == nil) && (parentDefinition.TargetNodeTypeName != nil) {
			self.TargetNodeTypeName = parentDefinition.TargetNodeTypeName
		}
		if (self.RelationshipDefinition == nil) && (parentDefinition.RelationshipDefinition != nil) {
			self.RelationshipDefinition = parentDefinition.RelationshipDefinition
		}
		if (self.Occurrences == nil) && (parentDefinition.Occurrences != nil) {
			self.Occurrences = parentDefinition.Occurrences
		}
		if (self.TargetCapabilityType == nil) && (parentDefinition.TargetCapabilityType != nil) {
			self.TargetCapabilityType = parentDefinition.TargetCapabilityType
		}
		if (self.TargetNodeType == nil) && (parentDefinition.TargetNodeType != nil) {
			self.TargetNodeType = parentDefinition.TargetNodeType
		}

		// Validate type compatibility
		if (self.TargetCapabilityType != nil) && (parentDefinition.TargetCapabilityType != nil) && !self.Context.Hierarchy.IsCompatible(parentDefinition.TargetCapabilityType, self.TargetCapabilityType) {
			self.Context.ReportIncompatibleType(self.TargetCapabilityType.Name, parentDefinition.TargetCapabilityType.Name)
		}
		if (self.TargetNodeType != nil) && (parentDefinition.TargetNodeType != nil) && !self.Context.Hierarchy.IsCompatible(parentDefinition.TargetNodeType, self.TargetNodeType) {
			self.Context.ReportIncompatibleType(self.TargetNodeType.Name, parentDefinition.TargetNodeType.Name)
		}
	}

	if self.TargetCapabilityTypeName == nil {
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
