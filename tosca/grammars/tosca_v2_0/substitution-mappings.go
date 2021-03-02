package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// SubstitutionMappings
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 2.10, 2.11, 2.12
// [TOSCA-Simple-Profile-YAML-v1.2] @ 2.10, 2.11
// [TOSCA-Simple-Profile-YAML-v1.1] @ 2.10, 2.11
// [TOSCA-Simple-Profile-YAML-v1.0] @ 2.10, 2.11
//

type SubstitutionMappings struct {
	*Entity `name:"substitution mappings"`

	NodeTypeName        *string             `read:"node_type" require:""`
	CapabilityMappings  CapabilityMappings  `read:"capabilities,CapabilityMapping"`
	RequirementMappings RequirementMappings `read:"requirements,RequirementMapping"`
	PropertyMappings    PropertyMappings    `read:"properties,PropertyMapping"`     // introduced in TOSCA 1.2
	AttributeMappings   AttributeMappings   `read:"attributes,AttributeMapping"`    // introduced in TOSCA 1.3
	InterfaceMappings   InterfaceMappings   `read:"interfaces,InterfaceMapping"`    // introduced in TOSCA 1.2
	SubstitutionFilter  NodeFilter          `read:"substitution_filter,NodeFilter"` // introduced in TOSCA 1.3

	NodeType *NodeType `lookup:"node_type,NodeTypeName" json:"-" yaml:"-"`
}

func NewSubstitutionMappings(context *tosca.Context) *SubstitutionMappings {
	return &SubstitutionMappings{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadSubstitutionMappings(context *tosca.Context) tosca.EntityPtr {
	if context.HasQuirk(tosca.QuirkSubstitutionMappingsRequirementsList) {
		context.SetReadTag("RequirementMappings", "requirements,{}RequirementMapping")
	}

	self := NewSubstitutionMappings(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

func (self *SubstitutionMappings) IsRequirementMapped(nodeTemplate *NodeTemplate, requirementName string) bool {
	for _, mapping := range self.RequirementMappings {
		if mapping.NodeTemplate == nodeTemplate {
			if (mapping.RequirementName != nil) && (*mapping.RequirementName == requirementName) {
				return true
			}
		}
	}
	return false
}

func (self *SubstitutionMappings) Normalize(normalServiceTemplate *normal.ServiceTemplate) *normal.Substitution {
	logNormalize.Debug("substitution mappings")

	if self.NodeType == nil {
		return nil
	}

	normalSubstitution := normalServiceTemplate.NewSubstitution()

	normalSubstitution.Type = tosca.GetCanonicalName(self.NodeType)

	if metadata, ok := self.NodeType.GetMetadata(); ok {
		normalSubstitution.TypeMetadata = metadata
	}

	for _, mapping := range self.CapabilityMappings {
		if (mapping.NodeTemplate == nil) || (mapping.CapabilityName == nil) {
			continue
		}

		if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]; ok {
			if normalCapability, ok := normalNodeTemplate.Capabilities[*mapping.CapabilityName]; ok {
				normalSubstitution.CapabilityMappings[normalNodeTemplate] = normalCapability
			}
		}
	}

	for _, mapping := range self.RequirementMappings {
		if (mapping.NodeTemplate == nil) || (mapping.RequirementName == nil) {
			continue
		}

		if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]; ok {
			normalSubstitution.RequirementMappings[normalNodeTemplate] = *mapping.RequirementName
		}
	}

	for _, mapping := range self.PropertyMappings {
		if (mapping.NodeTemplate == nil) || (mapping.PropertyName == nil) {
			continue
		}

		if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]; ok {
			normalServiceTemplate.Substitution.PropertyMappings[normalNodeTemplate] = *mapping.PropertyName
		}
	}

	for _, mapping := range self.AttributeMappings {
		if (mapping.NodeTemplate == nil) || (mapping.AttributeName == nil) {
			continue
		}

		if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]; ok {
			normalServiceTemplate.Substitution.AttributeMappings[normalNodeTemplate] = *mapping.AttributeName
		}
	}

	for _, mapping := range self.InterfaceMappings {
		if (mapping.NodeTemplate == nil) || (mapping.InterfaceName == nil) {
			continue
		}

		if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]; ok {
			normalServiceTemplate.Substitution.InterfaceMappings[normalNodeTemplate] = *mapping.InterfaceName
		}
	}

	return normalSubstitution
}
