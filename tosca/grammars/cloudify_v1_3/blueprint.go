package cloudify_v1_3

import (
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

	Groups []*Group `read:"groups,Group"`
}

func NewBlueprint(context *tosca.Context) *Blueprint {
	self := Blueprint{Unit: NewUnit(context)}

	self.Context.ImportScript("tosca.resolve", "internal:/tosca/simple/1.2/js/resolve.js")
	self.Context.ImportScript("tosca.coerce", "internal:/tosca/simple/1.2/js/coerce.js")
	self.Context.ImportScript("tosca.visualize", "internal:/tosca/simple/1.2/js/visualize.js")
	self.Context.ImportScript("tosca.utils", "internal:/tosca/simple/1.2/js/utils.js")
	self.Context.ImportScript("tosca.helpers", "internal:/tosca/simple/1.2/js/helpers.js")

	return &self
}

// tosca.Reader signature
func ReadBlueprint(context *tosca.Context) interface{} {
	self := NewBlueprint(context)
	context.ScriptNamespace.Merge(DefaultScriptNamespace)
	context.ValidateUnsupportedFields(append(context.ReadFields(self), "dsl_definitions"))
	return self
}

// tosca.Normalizable interface
func (self *Blueprint) Normalize() *normal.ServiceTemplate {
	log.Info("{normalize} template")

	s := normal.NewServiceTemplate()

	s.ScriptNamespace = self.Context.ScriptNamespace

	//self.Inputs.Normalize(s.Inputs, self.Context.FieldChild("inputs", nil))
	//self.Outputs.Normalize(s.Outputs, self.Context.FieldChild("outputs", nil))

	for _, nodeTemplate := range self.NodeTemplates {
		s.NodeTemplates[nodeTemplate.Name] = nodeTemplate.Normalize(s)
	}

	for _, nodeTemplate := range self.NodeTemplates {
		nodeTemplate.NormalizeRelationships(s)
	}

	return s
}
