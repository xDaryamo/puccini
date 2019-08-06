package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// OperationAssignment
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.15
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.13
//

type OperationAssignment struct {
	*Entity `name:"operation"`
	Name    string

	Description    *string                  `read:"description"`
	Implementation *OperationImplementation `read:"implementation,OperationImplementation"`
	Inputs         Values                   `read:"inputs,Value"`
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
		self.Implementation = ReadOperationImplementation(context.FieldChild("implementation", context.Data)).(*OperationImplementation)
	}

	return self
}

// tosca.Mappable interface
func (self *OperationAssignment) GetKey() string {
	return self.Name
}

func (self *OperationAssignment) Normalize(i *normal.Interface) *normal.Operation {
	log.Debugf("{normalize} operation: %s", self.Name)

	o := i.NewOperation(self.Name)

	if self.Description != nil {
		o.Description = *self.Description
	}

	if self.Implementation != nil {
		self.Implementation.Normalize(o)
	}

	self.Inputs.Normalize(o.Inputs)

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

		if assignment.Description == nil {
			assignment.Description = definition.Description
		}

		if (assignment.Implementation == nil) && (definition.Implementation != nil) {
			// If the definition has an implementation then we must have one, too
			assignment.Implementation = NewOperationImplementation(assignment.Context.FieldChild("implementation", nil))
		}

		if assignment.Implementation != nil {
			assignment.Implementation.Render(definition.Implementation)
		}

		assignment.Inputs.RenderProperties(definition.InputDefinitions, "input", assignment.Context.FieldChild("inputs", nil))
	}

	for key, assignment := range self {
		_, ok := definitions[key]
		if !ok {
			assignment.Context.ReportUndefined("operation")
			delete(self, key)
		}
	}
}

func (self OperationAssignments) Normalize(i *normal.Interface) {
	for key, operation := range self {
		i.Operations[key] = operation.Normalize(i)
	}
}
