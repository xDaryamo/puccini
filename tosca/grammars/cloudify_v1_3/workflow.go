package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// Workflow
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-workflows/]
//

type Workflow struct {
	*Entity `name:"workflow"`
	Name    string `namespace:""`

	Mapping    *string              `read:"mapping" require:"mapping"`
	Parameters ParameterDefinitions `read:"parameters,ParameterDefinition"`
}

func NewWorkflow(context *tosca.Context) *Workflow {
	return &Workflow{
		Entity:     NewEntity(context),
		Name:       context.Name,
		Parameters: make(ParameterDefinitions),
	}
}

// tosca.Reader signature
func ReadWorkflow(context *tosca.Context) interface{} {
	self := NewWorkflow(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}
