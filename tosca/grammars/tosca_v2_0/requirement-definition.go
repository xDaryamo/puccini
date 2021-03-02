package tosca_v2_0

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
)

//
// RequirementDefinition
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.2
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
func ReadRequirementDefinition(context *tosca.Context) tosca.EntityPtr {
	self := NewRequirementDefinition(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
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
	logInherit.Debugf("requirement definition: %s", self.Name)

	// Validate type compatibility
	if (self.TargetCapabilityType != nil) && (parentDefinition.TargetCapabilityType != nil) && !self.Context.Hierarchy.IsCompatible(parentDefinition.TargetCapabilityType, self.TargetCapabilityType) {
		self.Context.ReportIncompatibleType(self.TargetCapabilityType, parentDefinition.TargetCapabilityType)
	}
	if (self.TargetNodeType != nil) && (parentDefinition.TargetNodeType != nil) && !self.Context.Hierarchy.IsCompatible(parentDefinition.TargetNodeType, self.TargetNodeType) {
		self.Context.ReportIncompatibleType(self.TargetNodeType, parentDefinition.TargetNodeType)
	}

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

	if (self.RelationshipDefinition != nil) && (parentDefinition.RelationshipDefinition != nil) && (self.RelationshipDefinition != parentDefinition.RelationshipDefinition) {
		self.RelationshipDefinition.Inherit(parentDefinition.RelationshipDefinition)
	}
}

// parser.Renderable interface
func (self RequirementDefinition) Render() {
	logRender.Debugf("requirement definition: %s", self.Name)

	if self.TargetCapabilityTypeName == nil {
		// Avoid reporting more than once
		if !self.capabilityMissingProblemReported {
			self.Context.FieldChild("capability", nil).ReportFieldMissing()
			self.capabilityMissingProblemReported = true
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
		}
	}
}
