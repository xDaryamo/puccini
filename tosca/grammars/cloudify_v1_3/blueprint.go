package cloudify_v1_3

import (
	"fmt"

	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// Blueprint
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/]
//

type Blueprint struct {
	*Unit `name:"blueprint"`

	Description *string `read:"description"` // not in spec, but in code
	Groups      Groups  `read:"groups,Group"`
}

func NewBlueprint(context *tosca.Context) *Blueprint {
	return &Blueprint{Unit: NewUnit(context)}
}

// tosca.Reader signature
func ReadBlueprint(context *tosca.Context) interface{} {
	self := NewBlueprint(context)
	context.ScriptNamespace.Merge(DefaultScriptNamespace)
	context.ValidateUnsupportedFields(append(context.ReadFields(self), "dsl_definitions"))
	return self
}

// parser.HasInputs interface
func (self *Blueprint) SetInputs(inputs map[string]interface{}) {
	context := self.Context.FieldChild("inputs", nil)
	for name, data := range inputs {
		childContext := context.MapChild(name, data)
		if input, ok := self.Inputs[name]; ok {
			input.Value = ReadValue(childContext).(*Value)
		} else {
			childContext.ReportUndefined("input")
		}
	}
}

// tosca.Normalizable interface
func (self *Blueprint) Normalize() *normal.ServiceTemplate {
	log.Info("{normalize} blueprint")

	s := normal.NewServiceTemplate()

	if self.Metadata != nil {
		for k, v := range self.Metadata {
			// TODO: does Cloudify DSL really allow for any kind of value?
			s.Metadata[k] = fmt.Sprintf("%s", v)
		}
	}

	if self.Description != nil {
		s.Description = *self.Description
	}

	s.ScriptNamespace = self.Context.ScriptNamespace

	self.Inputs.Normalize(s.Inputs, self.Context.FieldChild("inputs", nil))
	self.Outputs.Normalize(s.Outputs)
	self.NodeTemplates.Normalize(s)
	self.Groups.Normalize(s)
	self.Workflows.Normalize(s)
	self.Policies.Normalize(s)

	// TODO: normalize plugins
	// TODO: normalize upload resources
	// TODO: normalize capabilities

	return s
}
