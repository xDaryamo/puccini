package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// OperationImplementation
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.14
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.13.2.3
//

type OperationImplementation struct {
	*Entity `name:"operation implementation"`

	Primary      *string   `read:"primary" require:"primary"`
	Dependencies *[]string `read:"dependencies"`
	// TODO: v1.2: timeout
	// TODO: v1.2: operation_host
}

func NewOperationImplementation(context *tosca.Context) *OperationImplementation {
	return &OperationImplementation{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadOperationImplementation(context *tosca.Context) interface{} {
	self := NewOperationImplementation(context)

	if context.Is("map") {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType("map", "string") {
		// Short notation
		self.Primary = context.ReadString()
	}

	return self
}

func (self *OperationImplementation) Normalize(o *normal.Operation) {
	if self.Primary != nil {
		o.Implementation = *self.Primary
	}

	if self.Dependencies != nil {
		o.Dependencies = *self.Dependencies
	}
}
