package v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/v1_1"
	"github.com/tliron/puccini/tosca/normal"
)

//
// SubstitutionMappings
//

type SubstitutionMappings struct {
	*v1_1.SubstitutionMappings

	PropertyMappings  []*PropertyMapping  `read:"properties,PropertyMapping"`
	InterfaceMappings []*InterfaceMapping `read:"interfaces,InterfaceMapping"`
}

func NewSubstitutionMappings(context *tosca.Context) *SubstitutionMappings {
	return &SubstitutionMappings{SubstitutionMappings: v1_1.NewSubstitutionMappings(context)}
}

// tosca.Reader signature
func ReadSubstitutionMappings(context *tosca.Context) interface{} {
	self := NewSubstitutionMappings(context)
	context.ValidateUnsupportedFields(context.ReadFields(self, Readers))
	return self
}

func (self *SubstitutionMappings) Normalize(s *normal.ServiceTemplate) {
	log.Info("{normalize} substitution mappings")

	if s.Substitution == nil {
		return
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
}
