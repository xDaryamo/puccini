package hot

import (
	"github.com/tliron/puccini/ard"
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
	Parameters           Parameters             `read:"parameters,Parameter"`
	Resources            []*Resource            `read:"resources,Resource"`
	Outputs              Outputs                `read:"outputs,Output"`
	ConditionDefinitions []*ConditionDefinition `read:"conditions,ConditionDefinition"`
}

func NewTemplate(context *tosca.Context) *Template {
	self := &Template{
		Entity:     NewEntity(context),
		Parameters: make(Parameters),
		Outputs:    make(Outputs),
	}

	self.Context.ImportScript("tosca.resolve", "internal:/tosca/simple/1.1/js/resolve.js")
	self.Context.ImportScript("tosca.coerce", "internal:/tosca/simple/1.1/js/coerce.js")
	self.Context.ImportScript("tosca.visualize", "internal:/tosca/simple/1.1/js/visualize.js")
	self.Context.ImportScript("tosca.utils", "internal:/tosca/simple/1.1/js/utils.js")
	self.Context.ImportScript("tosca.helpers", "internal:/tosca/simple/1.1/js/helpers.js")
	self.Context.ImportScript("openstack.generate", "internal:/tosca/openstack/1.0/js/generate.js")

	self.NewPseudoParameter("OS::stack_name", "stack_name")
	self.NewPseudoParameter("OS::stack_id", "stack_id")
	self.NewPseudoParameter("OS::project_id", "project_id")

	return self
}

// tosca.Reader signature
func ReadTemplate(context *tosca.Context) interface{} {
	self := NewTemplate(context)
	context.ScriptNamespace.Merge(DefaultScriptNamespace)
	context.ValidateUnsupportedFields(append(context.ReadFields(self, Readers)))
	return self
}

func (self *Template) NewPseudoParameter(name string, value string) {
	context := self.Context.FieldChild("parameters", nil).MapChild(name, ard.Map{
		"type":      "string",
		"hidden":    true,
		"immutable": true,
	})
	parameter := ReadParameter(context).(*Parameter)
	parameter.Value = NewValue(context.WithData(value))
	self.Parameters[name] = parameter
}

// parser.Importer interface
func (self *Template) GetImportSpecs() []*tosca.ImportSpec {
	var importSpecs []*tosca.ImportSpec
	return importSpecs
}

// parser.HasInputs interface
func (self *Template) SetInputs(inputs map[string]interface{}) {
	context := self.Context.FieldChild("parameters", nil)
	for name, data := range inputs {
		childContext := context.MapChild(name, data)
		parameter, ok := self.Parameters[name]
		if !ok {
			childContext.ReportUndefined("parameter")
			continue
		}

		parameter.Value = ReadValue(childContext).(*Value)
	}
}

// tosca.Normalizable interface
func (self *Template) Normalize() *normal.ServiceTemplate {
	log.Info("{normalize} template")

	s := normal.NewServiceTemplate()

	if self.Description != nil {
		s.Description = *self.Description
	}

	s.ScriptNamespace = self.Context.ScriptNamespace

	self.Parameters.Normalize(s.Inputs, self.Context.FieldChild("parameters", nil))
	self.Outputs.Normalize(s.Outputs, self.Context.FieldChild("outputs", nil))

	for _, resource := range self.Resources {
		s.NodeTemplates[resource.Name] = resource.Normalize(s)
	}

	for _, resource := range self.Resources {
		resource.NormalizeDependencies(s)
	}

	// TODO: normalize ParameterGroups

	return s
}
