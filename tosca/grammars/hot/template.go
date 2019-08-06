package hot

import (
	"time"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// Template
//
// [https://docs.openstack.org/heat/stein/template_guide/hot_spec.html#template-structure]
//

type Template struct {
	*Entity `name:"template"`

	HeatTemplateVersion  *string
	Description          *string              `read:"description"`
	ParameterGroups      ParameterGroups      `read:"parameter_groups,[]ParameterGroup"`
	Parameters           Parameters           `read:"parameters,Parameter"`
	Resources            Resources            `read:"resources,Resource"`
	Outputs              Outputs              `read:"outputs,Output"`
	ConditionDefinitions ConditionDefinitions `read:"conditions,ConditionDefinition"`
}

func NewTemplate(context *tosca.Context) *Template {
	self := &Template{
		Entity:     NewEntity(context),
		Parameters: make(Parameters),
		Outputs:    make(Outputs),
	}

	self.Context.ImportScript("tosca.resolve", "internal:/tosca/common/1.0/js/resolve.js")
	self.Context.ImportScript("tosca.coerce", "internal:/tosca/common/1.0/js/coerce.js")
	self.Context.ImportScript("tosca.utils", "internal:/tosca/common/1.0/js/utils.js")
	self.Context.ImportScript("tosca.helpers", "internal:/tosca/common/1.0/js/helpers.js")
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

	if heatTemplateVersionContext, ok := context.GetFieldChild("heat_template_version"); ok {
		switch heatTemplateVersionContext.Data.(type) {
		case time.Time:
			heatTemplateVersionContext.Data = heatTemplateVersionContext.Data.(time.Time).Format("2006-01-02")
		}

		if heatTemplateVersionContext.Is("string") {
			self.HeatTemplateVersion = heatTemplateVersionContext.ReadString()
		} else {
			heatTemplateVersionContext.ReportValueWrongType("string", "timestamp")
		}
	} else {
		context.FieldChild("heat_template_version", nil).ReportFieldMissing()
	}

	context.ValidateUnsupportedFields(append(context.ReadFields(self), "heat_template_version"))
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
		if parameter, ok := self.Parameters[name]; ok {
			parameter.Value = ReadValue(childContext).(*Value)
		} else {
			childContext.ReportUndefined("parameter")
		}
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
	self.Resources.Normalize(s)

	// TODO: normalize ParameterGroups

	return s
}
