package v1_1

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// OperationImplementation
//

type OperationImplementation struct {
	*Entity `name:"operation implementation"`

	Primary      *string   `read:"primary" require:"primary"`
	Dependencies *[]string `read:"dependencies"`
}

func NewOperationImplementation(context *tosca.Context) *OperationImplementation {
	return &OperationImplementation{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadOperationImplementation(context *tosca.Context) interface{} {
	self := NewOperationImplementation(context)
	if context.Is("map") {
		context.ValidateUnsupportedFields(context.ReadFields(self, Readers))
	} else if context.ValidateType("map", "string") {
		self.Primary = context.ReadString()
	}
	return self
}

func init() {
	Readers["OperationImplementation"] = ReadOperationImplementation
}

func (self *OperationImplementation) Normalize(o *normal.Operation) {
	if self.Primary != nil {
		o.Implementation = *self.Primary
	}

	if self.Dependencies != nil {
		o.Dependencies = *self.Dependencies
	}
}
