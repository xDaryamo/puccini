package v1_1

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// RequirementAssignment
//
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.2
//

type RequirementAssignment struct {
	*Entity `name:"requirement"`
	Name    string

	TargetNodeTemplateNameOrNodeTypeName     *string                 `read:"node"`
	TargetCapabilityNameOrCapabilityTypeName *string                 `read:"capability"`
	Relationship                             *RelationshipAssignment `read:"relationship,RelationshipAssignment"`
	TargetNodeFilter                         *NodeFilter             `read:"node_filter,NodeFilter"`

	TargetNodeTemplate   *NodeTemplate   `lookup:"node,TargetNodeTemplateNameOrNodeTypeName" json:"-" yaml:"-"`
	TargetNodeType       *NodeType       `lookup:"node,TargetNodeTemplateNameOrNodeTypeName" json:"-" yaml:"-"`
	TargetCapabilityType *CapabilityType `lookup:"capability,?TargetCapabilityNameOrCapabilityTypeName" json:"-" yaml:"-"`
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
		context.ValidateUnsupportedFields(context.ReadFields(self, Readers))
	} else if context.ValidateType("map", "string") {
		self.TargetNodeTemplateNameOrNodeTypeName = context.ReadString()
	}
	return self
}

func NewDefaultRequirementAssignment(definition *RequirementDefinition, context *tosca.Context) *RequirementAssignment {
	self := NewRequirementAssignment(context.MapChild(definition.Name, nil))
	self.TargetNodeTemplateNameOrNodeTypeName = definition.TargetNodeTypeName
	self.TargetNodeType = definition.TargetNodeType
	self.TargetCapabilityNameOrCapabilityTypeName = definition.TargetCapabilityTypeName
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

func (self *RequirementAssignment) Satisfy(s *normal.ServiceTemplate, n *normal.NodeTemplate, nodeTemplate *NodeTemplate, topologyTemplate *TopologyTemplate) {
	if (topologyTemplate.SubstitutionMappings != nil) && topologyTemplate.SubstitutionMappings.IsRequirementMapped(nodeTemplate, self.Name) {
		// Ignore mapped requirements
		log.Infof("{satisfy} %s: skipping because in substitution mappings", self.Context.Path)
		return
	}

	definition, ok := self.GetDefinition(nodeTemplate)
	if !ok {
		self.Context.ReportUnsatisfiedRequirement()
		return
	}

	// Candidate node templates

	var candidateNodeTemplates []*NodeTemplate

	if self.TargetNodeTemplate != nil {
		// Just this node template
		candidateNodeTemplates = []*NodeTemplate{self.TargetNodeTemplate}
	} else if self.TargetNodeType != nil {
		// Gather node templates of this type
		candidateNodeTemplates = topologyTemplate.GetNodeTemplatesOfType(self.TargetNodeType)
	} else if definition.TargetNodeType != nil {
		// Gather node templates of this type
		candidateNodeTemplates = topologyTemplate.GetNodeTemplatesOfType(definition.TargetNodeType)
	} else {
		// All node templates
		for _, nodeTemplate := range topologyTemplate.NodeTemplates {
			candidateNodeTemplates = append(candidateNodeTemplates, nodeTemplate)
		}
	}

	// Filter candidate node templates

	if nodeTemplate.NodeFilter != nil {
		candidateNodeTemplates = nodeTemplate.NodeFilter.FilterNodeTemplates(candidateNodeTemplates)
	}

	if self.TargetNodeFilter != nil {
		candidateNodeTemplates = self.TargetNodeFilter.FilterNodeTemplates(candidateNodeTemplates)
	}

	if len(candidateNodeTemplates) == 0 {
		log.Debugf("{satisfy} %s: no candidate node templates", self.Context.Path)
		self.Context.ReportUnsatisfiedRequirement()
		return
	}

	// Find first matching capability in candidate node templates

	for _, candidateNodeTemplate := range candidateNodeTemplates {
		if candidateNodeTemplate == nodeTemplate {
			// Don't satisfy requirements with self
			continue
		}

		log.Debugf("{satisfy} %s: trying node template: %s", self.Context.Path, candidateNodeTemplate.Name)

		var candidateCapabilities []*CapabilityAssignment

		if self.TargetCapabilityType != nil {
			// Gather capabilities of the specified type
			log.Debugf("{satisfy} %s: gathering \"%s\" capabilities in node template: %s", self.Context.Path, self.TargetCapabilityType.Name, candidateNodeTemplate.Name)
			candidateCapabilities = candidateNodeTemplate.GetCapabilitiesOfType(self.TargetCapabilityType)
		} else if self.TargetCapabilityNameOrCapabilityTypeName != nil {
			// Just this specified named capability
			if candidateCapability, ok := candidateNodeTemplate.GetCapabilityByName(*self.TargetCapabilityNameOrCapabilityTypeName, definition.TargetCapabilityType); ok {
				log.Debugf("{satisfy} %s: using capability named \"%s\" in node template: %s", self.Context.Path, candidateCapability.Name, candidateNodeTemplate.Name)
				candidateCapabilities = []*CapabilityAssignment{candidateCapability}
			} else {
				log.Debugf("{satisfy} %s: capability named %s is wrong type in node template: %s", self.Context.Path, self.TargetCapabilityNameOrCapabilityTypeName, candidateNodeTemplate.Name)
			}
		} else if definition.TargetCapabilityType != nil {
			// Gather capabilities of the type specified in the requirement definition
			log.Debugf("{satisfy} %s: gathering \"%s\" capabilities in node template: %s", self.Context.Path, definition.TargetCapabilityType.Name, candidateNodeTemplate.Name)
			candidateCapabilities = candidateNodeTemplate.GetCapabilitiesOfType(definition.TargetCapabilityType)
		}

		if len(candidateCapabilities) == 0 {
			log.Debugf("{satisfy} %s: no candidate capabilities in node template: %s", self.Context.Path, candidateNodeTemplate.Name)
			continue
		}

		// TODO: capability filter

		// TODO: check occurrences

		log.Infof("{satisfy} %s: satisfied", self.Context.Path)

		// Grab the first one
		capability := candidateCapabilities[0]
		r := n.NewRelationship()
		r.Name = self.Name
		r.TargetNodeTemplate = s.NodeTemplates[candidateNodeTemplate.Name]
		r.TargetCapability = r.TargetNodeTemplate.Capabilities[capability.Name]

		if self.Relationship != nil {
			self.Relationship.Normalize(r)
		}

		return
	}

	log.Infof("{satisfy} %s: not satisfied", self.Context.Path)

	self.Context.ReportUnsatisfiedRequirement()
}

//
// RequirementAssignments
//

type RequirementAssignments []*RequirementAssignment

func (self *RequirementAssignments) Render(definitions RequirementDefinitions, context *tosca.Context) {
	for key, definition := range definitions {
		found := false
		for _, assignment := range *self {
			if assignment.Name == key {
				found = true
				break
			}
		}

		if !found && (definition.Occurrences == nil) {
			// The TOSCA spec says that occurrences has an "implied default of [1,1]"
			// Our interpretation is that we should automatically add a single assignment if none was specified
			*self = append(*self, NewDefaultRequirementAssignment(definition, context))
		}
	}

	for index, assignment := range *self {
		if definition, ok := definitions[assignment.Name]; ok {
			if definition.RelationshipDefinition == nil {
				continue
			}

			if assignment.Relationship == nil {
				assignment.Relationship = definition.RelationshipDefinition.NewDefaultAssignment(assignment.Context.FieldChild("relationship", nil))
			}

			if assignment.Relationship.RelationshipTemplateOrRelationshipTypeName == nil {
				assignment.Relationship.RelationshipTemplateOrRelationshipTypeName = definition.RelationshipDefinition.RelationshipTypeName
			}

			if assignment.Relationship.RelationshipType == nil {
				assignment.Relationship.RelationshipType = definition.RelationshipDefinition.RelationshipType
			}

			assignment.Relationship.Render(definition.RelationshipDefinition)
		} else {
			assignment.Context.ReportUndefined("requirement")
			*self = append((*self)[:index], (*self)[index+1:]...)
		}
	}
}
