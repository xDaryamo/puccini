package cloudify_v1_3

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
)

//
// OperationDefinition
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-interfaces/]
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
func ReadOperationDefinition(context *tosca.Context) tosca.EntityPtr {
	self := NewOperationDefinition(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
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

func (self *OperationDefinition) Inherit(parentDefinition *OperationDefinition) {
	logInherit.Debugf("operation definition: %s", self.Name)

	self.InputDefinitions.Inherit(parentDefinition.InputDefinitions)
}

//
// OperationDefinitions
//

type OperationDefinitions map[string]*OperationDefinition

func (self OperationDefinitions) Inherit(parentDefinitions OperationDefinitions) {
	for name, definition := range parentDefinitions {
		if _, ok := self[name]; !ok {
			self[name] = definition
		}
	}

	for name, definition := range self {
		if parentDefinition, ok := parentDefinitions[name]; ok {
			if definition != parentDefinition {
				definition.Inherit(parentDefinition)
			}
		}
	}
}
