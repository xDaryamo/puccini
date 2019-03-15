package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// Input
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-inputs/]
//

type Input struct {
	*ParameterDefinition `name:"property definition"`
	Name                 string `namespace:""`
}

func NewInput(context *tosca.Context) *Input {
	return &Input{
		ParameterDefinition: NewParameterDefinition(context),
		Name:                context.Name,
	}
}

// tosca.Reader signature
func ReadInput(context *tosca.Context) interface{} {
	self := NewInput(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Mappable interface
func (self *Input) GetKey() string {
	return self.Name
}
