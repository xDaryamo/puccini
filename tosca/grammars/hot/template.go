package hot

import (
	"time"

	"github.com/tliron/kutil/ard"
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

	self.Context.ImportScriptlet("tosca.lib.utils", "internal:/tosca/common/1.0/js/lib/utils.js")
	self.Context.ImportScriptlet("tosca.lib.traversal", "internal:/tosca/common/1.0/js/lib/traversal.js")
	self.Context.ImportScriptlet("tosca.resolve", "internal:/tosca/common/1.0/js/resolve.js")
	self.Context.ImportScriptlet("tosca.coerce", "internal:/tosca/common/1.0/js/coerce.js")
	self.Context.ImportScriptlet("openstack.generate", "internal:/tosca/openstack/1.0/js/generate.js")

	self.NewPseudoParameter("OS::stack_name", "stack_name")
	self.NewPseudoParameter("OS::stack_id", "stack_id")
	self.NewPseudoParameter("OS::project_id", "project_id")

	return self
}

// tosca.Reader signature
func ReadTemplate(context *tosca.Context) tosca.EntityPtr {
	self := NewTemplate(context)
	context.ScriptletNamespace.Merge(DefaultScriptletNamespace)

	if heatTemplateVersionContext, ok := context.GetFieldChild("heat_template_version"); ok {
		switch data := heatTemplateVersionContext.Data.(type) {
		case time.Time:
			heatTemplateVersionContext.Data = data.Format("2006-01-02")
		}

		if heatTemplateVersionContext.Is(ard.TypeString) {
			self.HeatTemplateVersion = heatTemplateVersionContext.ReadString()
		} else {
			heatTemplateVersionContext.ReportValueWrongType(ard.TypeString, ard.TypeTimestamp)
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
	parameter.Value = NewValue(context.Clone(value))
	self.Parameters[name] = parameter
}

// parser.Importer interface
func (self *Template) GetImportSpecs() []*tosca.ImportSpec {
	var importSpecs []*tosca.ImportSpec
	return importSpecs
}

// parser.HasInputs interface
func (self *Template) SetInputs(inputs map[string]ard.Value) {
	context := self.Context.FieldChild("parameters", nil)
	for name, data := range inputs {
		childContext := context.MapChild(name, data)
		if parameter, ok := self.Parameters[name]; ok {
			parameter.Value = ReadValue(childContext).(*Value)
		} else {
			childContext.ReportUndeclared("parameter")
		}
	}
}

// normal.Normalizable interface
func (self *Template) NormalizeServiceTemplate() *normal.ServiceTemplate {
	logNormalize.Debug("template")

	normalServiceTemplate := normal.NewServiceTemplate()

	if self.Description != nil {
		normalServiceTemplate.Description = *self.Description
	}

	normalServiceTemplate.ScriptletNamespace = self.Context.ScriptletNamespace

	self.Parameters.Normalize(normalServiceTemplate.Inputs, self.Context.FieldChild("parameters", nil))
	self.Outputs.Normalize(normalServiceTemplate.Outputs, self.Context.FieldChild("outputs", nil))
	self.Resources.Normalize(normalServiceTemplate)

	// TODO: normalize ParameterGroups

	return normalServiceTemplate
}
