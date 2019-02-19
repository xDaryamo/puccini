package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_1"
	"github.com/tliron/puccini/tosca/normal"
)

//
// ServiceTemplate
//

type ServiceTemplate struct {
	*tosca_v1_1.ServiceTemplate

	TopologyTemplate *TopologyTemplate `read:"topology_template,TopologyTemplate"`
}

func NewServiceTemplate(context *tosca.Context) *ServiceTemplate {
	return &ServiceTemplate{ServiceTemplate: tosca_v1_1.NewServiceTemplate(context)}
}

// tosca.Reader signature
func ReadServiceTemplate(context *tosca.Context) interface{} {
	self := NewServiceTemplate(context)
	context.ScriptNamespace.Merge(tosca_v1_1.DefaultScriptNamespace)
	context.ValidateUnsupportedFields(append(context.ReadFields(self, Readers), "dsl_definitions"))

	// Hook up inheritance
	if self.TopologyTemplate != nil {
		self.ServiceTemplate.TopologyTemplate = self.TopologyTemplate.TopologyTemplate
	}

	return self
}

// tosca.Normalizable interface
func (self *ServiceTemplate) Normalize() *normal.ServiceTemplate {
	log.Info("{normalize} service template")

	s := self.ServiceTemplate.Normalize()

	// Hook up inheritance
	if (s != nil) && (self.TopologyTemplate != nil) && (self.TopologyTemplate.SubstitutionMappings != nil) {
		self.TopologyTemplate.SubstitutionMappings.Normalize(s)
	}

	return s
}
