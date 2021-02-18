package tosca_v2_0

import (
	"fmt"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// RequirementAssignment
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.2
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.2
//

type RequirementAssignment struct {
	*Entity `name:"requirement"`
	Name    string

	TargetCapabilityNameOrTypeName   *string                 `read:"capability"`
	TargetNodeTemplateNameOrTypeName *string                 `read:"node"`
	TargetNodeFilter                 *NodeFilter             `read:"node_filter,NodeFilter"`
	Relationship                     *RelationshipAssignment `read:"relationship,RelationshipAssignment"`
	Occurrences                      *RangeEntity            `read:"occurrences,RangeEntity"` // introduced in TOSCA 1.3

	TargetCapabilityType *CapabilityType `lookup:"capability,?TargetCapabilityNameOrTypeName" json:"-" yaml:"-"`
	TargetNodeTemplate   *NodeTemplate   `lookup:"node,TargetNodeTemplateNameOrTypeName" json:"-" yaml:"-"`
	TargetNodeType       *NodeType       `lookup:"node,TargetNodeTemplateNameOrTypeName" json:"-" yaml:"-"`
}

func NewRequirementAssignment(context *tosca.Context) *RequirementAssignment {
	return &RequirementAssignment{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadRequirementAssignment(context *tosca.Context) tosca.EntityPtr {
	self := NewRequirementAssignment(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.TargetNodeTemplateNameOrTypeName = context.FieldChild("node", context.Data).ReadString()
	}

	return self
}

func NewDefaultRequirementAssignment(index int, definition *RequirementDefinition, context *tosca.Context) *RequirementAssignment {
	context = context.SequencedListChild(index, definition.Name, nil)
	context.Name = definition.Name
	self := NewRequirementAssignment(context)
	self.TargetNodeTemplateNameOrTypeName = definition.TargetNodeTypeName
	self.TargetNodeType = definition.TargetNodeType
	self.TargetCapabilityNameOrTypeName = definition.TargetCapabilityTypeName
	self.TargetCapabilityType = definition.TargetCapabilityType
	return self
}

func (self *RequirementAssignment) GetDefinition(nodeTemplate *NodeTemplate) (*RequirementDefinition, bool) {
	if nodeTemplate.NodeType == nil {
		return nil, false
	}
	definition, ok := nodeTemplate.NodeType.RequirementDefinitions[self.Name]
	return definition, ok
}

func (self *RequirementAssignment) Normalize(nodeTemplate *NodeTemplate, normalNodeTemplate *normal.NodeTemplate) *normal.Requirement {
	normalRequirement := normalNodeTemplate.NewRequirement(self.Name, normal.NewLocationForContext(self.Context))

	if self.TargetCapabilityType != nil {
		name := tosca.GetCanonicalName(self.TargetCapabilityType)
		normalRequirement.CapabilityTypeName = &name
	} else if self.TargetCapabilityNameOrTypeName != nil {
		normalRequirement.CapabilityName = self.TargetCapabilityNameOrTypeName
	}

	if self.TargetNodeTemplate != nil {
		normalRequirement.NodeTemplate, _ = normalNodeTemplate.ServiceTemplate.NodeTemplates[self.TargetNodeTemplate.Name]
	}

	if self.TargetNodeType != nil {
		name := tosca.GetCanonicalName(self.TargetNodeType)
		normalRequirement.NodeTypeName = &name
	}

	if nodeTemplate.RequirementTargetsNodeFilter != nil {
		nodeTemplate.RequirementTargetsNodeFilter.Normalize(normalRequirement)
	}

	if self.TargetNodeFilter != nil {
		self.TargetNodeFilter.Normalize(normalRequirement)
	}

	if self.Relationship != nil {
		if definition, ok := self.GetDefinition(nodeTemplate); ok {
			self.Relationship.Normalize(definition.RelationshipDefinition, normalRequirement.NewRelationship())
		} else {
			self.Relationship.Normalize(nil, normalRequirement.NewRelationship())
		}
	}

	return normalRequirement
}

//
// RequirementAssignments
//

type RequirementAssignments []*RequirementAssignment

func (self *RequirementAssignments) Render(definitions RequirementDefinitions, context *tosca.Context) {
	// TODO: currently have no idea what to do with "occurrences" keyword in the requirement
	// assignment, because we interpret "occurrences" in the definition to mean how many times
	// it would be assigned

	for key, definition := range definitions {
		if definition.Occurrences == nil {
			// The TOSCA spec says that occurrences has an "implied default of [1,1]"
			// Our interpretation is that we should automatically add a single assignment if none was specified

			found := false
			for _, assignment := range *self {
				if assignment.Name == key {
					found = true
					break
				}
			}

			if !found {
				*self = append(*self, NewDefaultRequirementAssignment(len(*self), definition, context))
			}
		} else {
			// Check occurrences
			var count uint64
			for _, assignment := range *self {
				if assignment.Name == key {
					count++
				}
			}

			if !definition.Occurrences.Range.InRange(count) {
				context.ReportNotInRange(fmt.Sprintf("number of requirement %q assignments", definition.Name), count, definition.Occurrences.Range.Lower, definition.Occurrences.Range.Upper)
			}
		}
	}

	for _, assignment := range *self {
		if definition, ok := definitions[assignment.Name]; ok {
			if assignment.TargetCapabilityNameOrTypeName == nil {
				assignment.TargetCapabilityNameOrTypeName = definition.TargetCapabilityTypeName
			}

			if assignment.TargetCapabilityType == nil {
				assignment.TargetCapabilityType = definition.TargetCapabilityType
			}

			if assignment.TargetNodeTemplateNameOrTypeName == nil {
				assignment.TargetNodeTemplateNameOrTypeName = definition.TargetNodeTypeName
			}

			if assignment.TargetNodeType == nil {
				assignment.TargetNodeType = definition.TargetNodeType
			}

			if definition.RelationshipDefinition != nil {
				if assignment.Relationship == nil {
					assignment.Relationship = definition.RelationshipDefinition.NewDefaultAssignment(assignment.Context.FieldChild("relationship", nil))
				}

				if assignment.Relationship.RelationshipTemplateNameOrTypeName == nil {
					// Note: the definition can only specify a relationship type, not a relationship template
					assignment.Relationship.RelationshipTemplateNameOrTypeName = definition.RelationshipDefinition.RelationshipTypeName
				}

				if (assignment.Relationship.RelationshipType == nil) && (assignment.Relationship.RelationshipTemplate == nil) {
					// Note: we are careful not set the relationship type if the assignment uses a relationship template
					assignment.Relationship.RelationshipType = definition.RelationshipDefinition.RelationshipType
				}
			}

			if assignment.Relationship != nil {
				assignment.Relationship.Render(definition.RelationshipDefinition)
			}
		} else {
			assignment.Context.ReportUndeclared("requirement")
		}
	}
}

func (self RequirementAssignments) Normalize(nodeTemplate *NodeTemplate, normalNodeTemplate *normal.NodeTemplate) {
	for _, requirement := range self {
		requirement.Normalize(nodeTemplate, normalNodeTemplate)
	}
}
