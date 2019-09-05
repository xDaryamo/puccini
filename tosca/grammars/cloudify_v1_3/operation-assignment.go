package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// OperationAssignment
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-interfaces/]
//

type OperationAssignment struct {
	*Entity `name:"operation assignment"`
	Name    string

	Implementation *string  `read:"implementation" require:"implementation"`
	Inputs         Values   `read:"inputs,Value"`
	Executor       *string  `read:"executor"`
	MaxRetries     *int64   `read:"max_retries"`
	RetryInterval  *float64 `read:"retry_interval"`
}

func NewOperationAssignment(context *tosca.Context) *OperationAssignment {
	return &OperationAssignment{
		Entity: NewEntity(context),
		Name:   context.Name,
		Inputs: make(Values),
	}
}

// tosca.Reader signature
func ReadOperationAssignment(context *tosca.Context) interface{} {
	self := NewOperationAssignment(context)

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

func ValidateOperationExecutor(executor string, context *tosca.Context) {
	switch executor {
	case "central_deployment_agent", "host_agent":
	default:
		context.FieldChild("executor", executor).ReportFieldUnsupportedValue()
	}
}

// tosca.Mappable interface
func (self *OperationAssignment) GetKey() string {
	return self.Name
}

func (self *OperationAssignment) Normalize(i *normal.Interface) *normal.Operation {
	log.Debugf("{normalize} operation: %s", self.Name)

	o := i.NewOperation(self.Name)

	if self.Implementation != nil {
		o.Implementation = *self.Implementation
	}

	self.Inputs.Normalize(o.Inputs, "")

	return o
}

//
// OperationAssignments
//

type OperationAssignments map[string]*OperationAssignment

func (self OperationAssignments) Render(definitions OperationDefinitions, context *tosca.Context) {
	for key, definition := range definitions {
		assignment, ok := self[key]

		if !ok {
			assignment = NewOperationAssignment(context.FieldChild(key, nil))
			self[key] = assignment
		}

		if (assignment.Implementation == nil) && (definition.Implementation != nil) {
			assignment.Implementation = definition.Implementation
		}

		assignment.Inputs.RenderParameters(definition.InputDefinitions, "input", assignment.Context.FieldChild("inputs", nil))
	}

	for key, assignment := range self {
		_, ok := definitions[key]
		if !ok {
			assignment.Context.ReportUndeclared("operation")
			delete(self, key)
		}
	}
}

func (self OperationAssignments) Normalize(i *normal.Interface) {
	for key, operation := range self {
		i.Operations[key] = operation.Normalize(i)
	}
}
