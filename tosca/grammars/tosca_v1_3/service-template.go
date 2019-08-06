package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// ServiceTemplate
//
// See Unit
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.10
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.9
//

type ServiceTemplate struct {
	*Unit `name:"service template"`

	Description      *string           `read:"description"`
	TopologyTemplate *TopologyTemplate `read:"topology_template,TopologyTemplate"`
}

func NewServiceTemplate(context *tosca.Context) *ServiceTemplate {
	return &ServiceTemplate{Unit: NewUnit(context)}
}

// tosca.Reader signature
func ReadServiceTemplate(context *tosca.Context) interface{} {
	self := NewServiceTemplate(context)
	context.ScriptNamespace.Merge(DefaultScriptNamespace)
	context.ValidateUnsupportedFields(append(context.ReadFields(self), "dsl_definitions"))
	return self
}

// tosca.Normalizable interface
func (self *ServiceTemplate) Normalize() *normal.ServiceTemplate {
	log.Info("{normalize} service template")

	s := normal.NewServiceTemplate()

	if self.Description != nil {
		s.Description = *self.Description
	}

	s.ScriptNamespace = self.Context.ScriptNamespace

	self.Unit.Normalize(s)
	if self.TopologyTemplate != nil {
		self.TopologyTemplate.Normalize(s)
	}

	return s
}
