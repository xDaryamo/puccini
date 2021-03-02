package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// ServiceTemplate
//
// See Unit
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.10
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.10
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.9
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.9
//

type ServiceTemplate struct {
	*Unit `name:"service template"`

	TopologyTemplate *TopologyTemplate `read:"topology_template,TopologyTemplate"`
}

func NewServiceTemplate(context *tosca.Context) *ServiceTemplate {
	return &ServiceTemplate{Unit: NewUnit(context)}
}

// tosca.Reader signature
func ReadServiceTemplate(context *tosca.Context) tosca.EntityPtr {
	self := NewServiceTemplate(context)
	context.ScriptletNamespace.Merge(DefaultScriptletNamespace)
	context.ValidateUnsupportedFields(append(context.ReadFields(self), "dsl_definitions"))
	if self.Profile != nil {
		context.CanonicalNamespace = self.Profile
	}
	return self
}

// normal.Normalizable interface
func (self *ServiceTemplate) NormalizeServiceTemplate() *normal.ServiceTemplate {
	logNormalize.Debug("service template")

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
