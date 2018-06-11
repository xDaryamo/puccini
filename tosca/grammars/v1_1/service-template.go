package v1_1

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// ServiceTemplate
//

type ServiceTemplate struct {
	*Unit `name:"service template"`

	TopologyTemplate *TopologyTemplate `read:"topology_template,TopologyTemplate"`
}

func NewServiceTemplate(context *tosca.Context) *ServiceTemplate {
	return &ServiceTemplate{Unit: NewUnit(context)}
}

// tosca.Reader signature
func ReadServiceTemplate(context *tosca.Context) interface{} {
	self := NewServiceTemplate(context)
	context.ScriptNamespace.Merge(DefaultScriptNamespace)
	context.ValidateUnsupportedFields(append(context.ReadFields(self, Readers), "dsl_definitions"))
	return self
}

func init() {
	Readers["ServiceTemplate"] = ReadServiceTemplate
}

// tosca.Normalizable interface
func (self *ServiceTemplate) Normalize() *normal.ServiceTemplate {
	log.Info("{normalize} service template")

	s := normal.NewServiceTemplate()

	s.ScriptNamespace = self.Context.ScriptNamespace

	self.Unit.Normalize(s)
	if self.TopologyTemplate != nil {
		self.TopologyTemplate.Normalize(s)
	}

	return s
}
