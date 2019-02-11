package v2018_08_31

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// Template
//
// [https://docs.openstack.org/heat/rocky/template_guide/hot_spec.html#template-structure]
//

type Template struct {
	*Entity `name:"template"`

	HeatTemplateVersion *string           `read:"heat_template_version" require:"heat_template_version"`
	Description         *string           `read:"description"`
	ParameterGroups     []*ParameterGroup `read:"parameter_groups,ParameterGroup"`
}

func NewTemplate(context *tosca.Context) *Template {
	return &Template{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadTemplate(context *tosca.Context) interface{} {
	self := NewTemplate(context)
	context.ScriptNamespace.Merge(DefaultScriptNamespace)
	context.ValidateUnsupportedFields(append(context.ReadFields(self, Readers)))
	return self
}

// tosca.Normalizable interface
func (self *Template) Normalize() *normal.ServiceTemplate {
	log.Info("{normalize} template")

	s := normal.NewServiceTemplate()

	if self.Description != nil {
		s.Description = *self.Description
	}

	s.ScriptNamespace = self.Context.ScriptNamespace

	return s
}
