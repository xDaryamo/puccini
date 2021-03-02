package cloudify_v1_3

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// OperationAssignment
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-interfaces/]
//

type OperationAssignment struct {
	*Entity `name:"operation assignment"`
	Name    string

	Implementation *string  `read:"implementation" require:""`
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
func ReadOperationAssignment(context *tosca.Context) tosca.EntityPtr {
	self := NewOperationAssignment(context)

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

func (self *OperationAssignment) Normalize(normalInterface *normal.Interface) *normal.Operation {
	logNormalize.Debugf("operation: %s", self.Name)

	normalOperation := normalInterface.NewOperation(self.Name)

	if self.Implementation != nil {
		normalOperation.Implementation = *self.Implementation
	}

	self.Inputs.Normalize(normalOperation.Inputs, "")

	return normalOperation
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

func (self OperationAssignments) Normalize(normalInterface *normal.Interface) {
	for key, operation := range self {
		normalInterface.Operations[key] = operation.Normalize(normalInterface)
	}
}
