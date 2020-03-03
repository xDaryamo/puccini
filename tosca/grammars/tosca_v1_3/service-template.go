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
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.10
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.10
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.9
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.9
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
	context.ScriptletNamespace.Merge(DefaultScriptletNamespace)
	context.ValidateUnsupportedFields(append(context.ReadFields(self), "dsl_definitions"))
	return self
}

// normal.Normalizable interface
func (self *ServiceTemplate) NormalizeServiceTemplate() *normal.ServiceTemplate {
	log.Info("{normalize} service template")

	normalServiceTemplate := normal.NewServiceTemplate()

	if self.Description != nil {
		normalServiceTemplate.Description = *self.Description
	}

	normalServiceTemplate.ScriptletNamespace = self.Context.ScriptletNamespace

	self.Unit.Normalize(normalServiceTemplate)
	if self.TopologyTemplate != nil {
		self.TopologyTemplate.Normalize(normalServiceTemplate)
	}

	return normalServiceTemplate
}
