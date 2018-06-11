package v1_1

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// SubstitutionMappings
//

type SubstitutionMappings struct {
	*Entity `name:"substitution mappings"`

	NodeTypeName        *string               `read:"node_type" require:"node_type"`
	CapabilityMappings  []*CapabilityMapping  `read:"capabilities,CapabilityMapping"`
	RequirementMappings []*RequirementMapping `read:"requirements,RequirementMapping"`

	NodeType *NodeType `lookup:"node_type,NodeTypeName" json:"-" yaml:"-"`
}

func NewSubstitutionMappings(context *tosca.Context) *SubstitutionMappings {
	return &SubstitutionMappings{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadSubstitutionMappings(context *tosca.Context) interface{} {
	self := NewSubstitutionMappings(context)
	context.ValidateUnsupportedFields(context.ReadFields(self, Readers))
	return self
}

func init() {
	Readers["SubstitutionMappings"] = ReadSubstitutionMappings
}

func (self *SubstitutionMappings) IsRequirementMapped(nodeTemplate *NodeTemplate, requirementName string) bool {
	for _, mappedRequirement := range self.RequirementMappings {
		if mappedRequirement.NodeTemplate == nodeTemplate {
			if (mappedRequirement.RequirementName != nil) && (*mappedRequirement.RequirementName == requirementName) {
				return true
			}
		}
	}
	return false
}

func (self *SubstitutionMappings) Normalize(s *normal.ServiceTemplate) {
	log.Info("{normalize} substitution mappings")

	if self.NodeType == nil {
		return
	}

	t := s.NewSubstitution()

	t.Type = self.NodeType.Name

	if metadata, ok := self.NodeType.GetMetadata(); ok {
		t.TypeMetadata = metadata
	}

	for _, mapped := range self.CapabilityMappings {
		if (mapped.NodeTemplate == nil) || (mapped.CapabilityName == nil) {
			continue
		}

		if n, ok := s.NodeTemplates[mapped.NodeTemplate.Name]; ok {
			if c, ok := n.Capabilities[*mapped.CapabilityName]; ok {
				t.CapabilityMappings[n] = c
			}
		}
	}

	for _, mapped := range self.RequirementMappings {
		if (mapped.NodeTemplate == nil) || (mapped.RequirementName == nil) {
			continue
		}

		if n, ok := s.NodeTemplates[mapped.NodeTemplate.Name]; ok {
			t.RequirementMappings[n] = *mapped.RequirementName
		}
	}
}
