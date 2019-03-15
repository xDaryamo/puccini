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

	Implementation *string `read:"implementation" require:"implementation"`
	Inputs         Values  `read:"properties,Value"`
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
		self.Implementation = context.ReadString()
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

	if self.Implementation != nil {
		o.Implementation = *self.Implementation
	}

	self.Inputs.Normalize(o.Inputs)

	return o
}

//
// OperationAssignments
//

type OperationAssignments map[string]*OperationAssignment
