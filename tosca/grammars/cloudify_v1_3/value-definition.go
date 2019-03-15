package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// ValueDefinition
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-capabilities/]
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-outputs/]
//

type ValueDefinition struct {
	*Entity `name:"capability"`
	Name    string `namespace:""`

	Description *string `read:"description"`
	Value       *Value  `read:"value,Value"`
}

func NewValueDefinition(context *tosca.Context) *ValueDefinition {
	return &ValueDefinition{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadValueDefinition(context *tosca.Context) interface{} {
	self := NewValueDefinition(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Mappable interface
func (self *ValueDefinition) GetKey() string {
	return self.Name
}
