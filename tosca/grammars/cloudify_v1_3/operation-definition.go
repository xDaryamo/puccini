package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// OperationDefinition
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-interfaces/]
//

type OperationDefinition struct {
	*Entity `name:"operation definition"`
	Name    string

	Implementation   *string              `read:"implementation"`
	InputDefinitions ParameterDefinitions `read:"inputs,ParameterDefinition"`
	Executor         *string              `read:"executor"`
	MaxRetries       *int64               `read:"max_retries"`
	RetryInterval    *float64             `read:"retry_interval"`
}

func NewOperationDefinition(context *tosca.Context) *OperationDefinition {
	return &OperationDefinition{
		Entity:           NewEntity(context),
		Name:             context.Name,
		InputDefinitions: make(ParameterDefinitions),
	}
}

// tosca.Reader signature
func ReadOperationDefinition(context *tosca.Context) interface{} {
	self := NewOperationDefinition(context)

	if context.Is("map") {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType("map", "string") {
		// Short notation
		self.Implementation = context.FieldChild("implementation", context.Data).ReadString()
	}

	if self.Executor != nil {
		ValidateOperationExecutor(*self.Executor, self.Context)
	}

	return self
}

// tosca.Mappable interface
func (self *OperationDefinition) GetKey() string {
	return self.Name
}

//
// OperationDefinitions
//

type OperationDefinitions map[string]*OperationDefinition
