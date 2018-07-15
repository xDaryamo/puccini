package v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/v1_1"
)

//
// TopologyTemplate
//

type TopologyTemplate struct {
	*v1_1.TopologyTemplate

	SubstitutionMappings *SubstitutionMappings `read:"substitution_mappings,SubstitutionMappings"`
}

func NewTopologyTemplate(context *tosca.Context) *TopologyTemplate {
	return &TopologyTemplate{TopologyTemplate: v1_1.NewTopologyTemplate(context)}
}

// tosca.Reader signature
func ReadTopologyTemplate(context *tosca.Context) interface{} {
	self := NewTopologyTemplate(context)
	context.ValidateUnsupportedFields(context.ReadFields(self, Readers))

	// Hook up inheritance
	if self.SubstitutionMappings != nil {
		self.TopologyTemplate.SubstitutionMappings = self.SubstitutionMappings.SubstitutionMappings
	}

	return self
}
