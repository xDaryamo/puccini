package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// SubstitutionMappings
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 2.10
// [TOSCA-Simple-Profile-YAML-v1.2] @ 2.11
// [TOSCA-Simple-Profile-YAML-v1.1] @ 2.10
// [TOSCA-Simple-Profile-YAML-v1.1] @ 2.11
//

type SubstitutionMappings struct {
	*Entity `name:"substitution mappings"`

	NodeTypeName        *string             `read:"node_type" require:"node_type"`
	CapabilityMappings  CapabilityMappings  `read:"capabilities,CapabilityMapping"`
	RequirementMappings RequirementMappings `read:"requirements,RequirementMapping"`
	PropertyMappings    PropertyMappings    `read:"properties,PropertyMapping"`
	InterfaceMappings   InterfaceMappings   `read:"interfaces,InterfaceMapping"`

	NodeType *NodeType `lookup:"node_type,NodeTypeName" json:"-" yaml:"-"`
}

func NewSubstitutionMappings(context *tosca.Context) *SubstitutionMappings {
	return &SubstitutionMappings{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadSubstitutionMappings(context *tosca.Context) interface{} {
	if context.HasQuirk("substitution_mappings.requirements.list") {
		if context.ReadOverrides == nil {
			context.ReadOverrides = make(map[string]string)
		}
		context.ReadOverrides["RequirementMappings"] = "requirements,{}RequirementMapping"
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

func (self *SubstitutionMappings) Normalize(s *normal.ServiceTemplate) *normal.Substitution {
	log.Info("{normalize} substitution mappings")

	if self.NodeType == nil {
		return nil
	}

	t := s.NewSubstitution()

	t.Type = self.NodeType.Name

	if metadata, ok := self.NodeType.GetMetadata(); ok {
		t.TypeMetadata = metadata
	}

	for _, mapping := range self.CapabilityMappings {
		if (mapping.NodeTemplate == nil) || (mapping.CapabilityName == nil) {
			continue
		}

		if n, ok := s.NodeTemplates[mapping.NodeTemplate.Name]; ok {
			if c, ok := n.Capabilities[*mapping.CapabilityName]; ok {
				t.CapabilityMappings[n] = c
			}
		}
	}

	for _, mapping := range self.RequirementMappings {
		if (mapping.NodeTemplate == nil) || (mapping.RequirementName == nil) {
			continue
		}

		if n, ok := s.NodeTemplates[mapping.NodeTemplate.Name]; ok {
			t.RequirementMappings[n] = *mapping.RequirementName
		}
	}

	for _, mapping := range self.PropertyMappings {
		if (mapping.NodeTemplate == nil) || (mapping.PropertyName == nil) {
			continue
		}

		if n, ok := s.NodeTemplates[mapping.NodeTemplate.Name]; ok {
			s.Substitution.PropertyMappings[n] = *mapping.PropertyName
		}
	}

	for _, mapping := range self.InterfaceMappings {
		if (mapping.NodeTemplate == nil) || (mapping.InterfaceName == nil) {
			continue
		}

		if n, ok := s.NodeTemplates[mapping.NodeTemplate.Name]; ok {
			s.Substitution.InterfaceMappings[n] = *mapping.InterfaceName
		}
	}

	return t
}
