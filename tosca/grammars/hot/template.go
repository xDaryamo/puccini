package hot

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

	HeatTemplateVersion  *string                `read:"heat_template_version" require:"heat_template_version"`
	Description          *string                `read:"description"`
	ParameterGroups      []*ParameterGroup      `read:"parameter_groups,[]ParameterGroup"`
	Parameters           []*Parameter           `read:"parameters,Parameter"`
	Resources            []*Resource            `read:"resources,Resource"`
	Outputs              []*Output              `read:"outputs,Output"`
	ConditionDefinitions []*ConditionDefinition `read:"conditions,ConditionDefinition"`
}

func NewTemplate(context *tosca.Context) *Template {
	return &Template{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadTemplate(context *tosca.Context) interface{} {
	self := NewTemplate(context)
	context.ImportScript("tosca.resolve", "internal:/tosca/simple/1.1/js/resolve.js")
	context.ImportScript("tosca.resolve", "internal:/tosca/simple/1.1/js/utils.js")
	context.ImportScript("tosca.resolve", "internal:/tosca/simple/1.1/js/helpers.js")
	context.ScriptNamespace.Merge(DefaultScriptNamespace)
	context.ValidateUnsupportedFields(append(context.ReadFields(self, Readers)))
	return self
}

// parser.Importer interface
func (self *Template) GetImportSpecs() []*tosca.ImportSpec {
	var importSpecs []*tosca.ImportSpec
	return importSpecs
}

// tosca.Normalizable interface
func (self *Template) Normalize() *normal.ServiceTemplate {
	log.Info("{normalize} template")

	s := normal.NewServiceTemplate()

	if self.Description != nil {
		s.Description = *self.Description
	}

	s.ScriptNamespace = self.Context.ScriptNamespace

	for _, resource := range self.Resources {
		s.NodeTemplates[resource.Name] = resource.Normalize(s)
	}

	return s
}
