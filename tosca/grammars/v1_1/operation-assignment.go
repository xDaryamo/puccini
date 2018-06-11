package v1_1

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// OperationAssignment
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
		context.ValidateUnsupportedFields(context.ReadFields(self, Readers))
	} else if context.ValidateType("map", "string") {
		self.Implementation = ReadOperationImplementation(context).(*OperationImplementation)
	}
	return self
}

func init() {
	Readers["OperationAssignment"] = ReadOperationAssignment
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
			assignment = NewOperationAssignment(context.MapChild(key, nil))
			self[key] = assignment
		}
		if assignment.Description == nil {
			assignment.Description = definition.Description
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
