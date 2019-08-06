package tosca_v1_3

import (
	"fmt"

	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// RequirementAssignment
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.2
//

type RequirementAssignment struct {
	*Entity `name:"requirement"`
	Name    string

	TargetCapabilityNameOrTypeName   *string                 `read:"capability"`
	TargetNodeTemplateNameOrTypeName *string                 `read:"node"`
	TargetNodeFilter                 *NodeFilter             `read:"node_filter,NodeFilter"`
	Relationship                     *RelationshipAssignment `read:"relationship,RelationshipAssignment"`

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
func ReadRequirementAssignment(context *tosca.Context) interface{} {
	self := NewRequirementAssignment(context)

	if context.Is("map") {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType("map", "string") {
		// Short notation
		self.TargetNodeTemplateNameOrTypeName = context.FieldChild("node", context.Data).ReadString()
	}

	return self
}

func NewDefaultRequirementAssignment(index int, definition *RequirementDefinition, context *tosca.Context) *RequirementAssignment {
	context = context.ListChild(index, nil)
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

func (self *RequirementAssignment) Normalize(nodeTemplate *NodeTemplate, s *normal.ServiceTemplate, n *normal.NodeTemplate) *normal.Requirement {
	r := n.NewRequirement(self.Name, self.Context.Path.String())

	if self.TargetCapabilityType != nil {
		r.CapabilityTypeName = &self.TargetCapabilityType.Name
	} else if self.TargetCapabilityNameOrTypeName != nil {
		r.CapabilityName = self.TargetCapabilityNameOrTypeName
	}

	if self.TargetNodeType != nil {
		r.NodeTypeName = &self.TargetNodeType.Name
	} else if self.TargetNodeTemplate != nil {
		r.NodeTemplate = s.NodeTemplates[self.TargetNodeTemplate.Name]
	}

	if nodeTemplate.RequirementTargetsNodeFilter != nil {
		nodeTemplate.RequirementTargetsNodeFilter.Normalize(r)
	}

	if self.TargetNodeFilter != nil {
		self.TargetNodeFilter.Normalize(r)
	}

	if self.Relationship != nil {
		self.Relationship.Normalize(r.NewRelationship())
	}

	return r
}

//
// RequirementAssignments
//

type RequirementAssignments []*RequirementAssignment

func (self *RequirementAssignments) Render(definitions RequirementDefinitions, context *tosca.Context) {
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
				context.ReportNotInRange(fmt.Sprintf("number of requirement \"%s\" assignments", definition.Name), count, definition.Occurrences.Range.Lower, definition.Occurrences.Range.Upper)
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
					assignment.Relationship.RelationshipTemplateNameOrTypeName = definition.RelationshipDefinition.RelationshipTypeName
				}

				if assignment.Relationship.RelationshipType == nil {
					assignment.Relationship.RelationshipType = definition.RelationshipDefinition.RelationshipType
				}

				assignment.Relationship.Render(definition.RelationshipDefinition)
			}
		} else {
			assignment.Context.ReportUndefined("requirement")
			// TODO: move to outside of loop?
			//*self = append((*self)[:index], (*self)[index+1:]...)
		}
	}
}

func (self RequirementAssignments) Normalize(nodeTemplate *NodeTemplate, s *normal.ServiceTemplate, n *normal.NodeTemplate) {
	for _, requirement := range self {
		requirement.Normalize(nodeTemplate, s, n)
	}
}
