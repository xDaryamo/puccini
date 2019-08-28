package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// InterfaceImplementation
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.14
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.13.2.3
//

type InterfaceImplementation struct {
	*Entity `name:"operation implementation"`

	Primary       *string   `read:"primary"`
	Dependencies  *[]string `read:"dependencies"`
	Timeout       *int64    `read:"timeout"`
	OperationHost *string   `read:"operation_host"`
}

func NewInterfaceImplementation(context *tosca.Context) *InterfaceImplementation {
	return &InterfaceImplementation{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadInterfaceImplementation(context *tosca.Context) interface{} {
	self := NewInterfaceImplementation(context)

	if context.Is("map") {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType("map", "string") {
		// Short notation
		self.Primary = context.FieldChild("primary", context.Data).ReadString()
	}

	return self
}

func (self *InterfaceImplementation) Render(definition *InterfaceImplementation) {
	if definition != nil {
		if (self.Primary != nil) && (definition.Primary != nil) {
			self.Primary = definition.Primary
		}
		if (self.Dependencies != nil) && (definition.Dependencies != nil) {
			self.Dependencies = definition.Dependencies
		}
		if (self.Timeout != nil) && (definition.Timeout != nil) {
			self.Timeout = definition.Timeout
		}
		if (self.OperationHost != nil) && (definition.OperationHost != nil) {
			self.OperationHost = definition.OperationHost
		}
	}

	if self.OperationHost != nil {
		path := self.Context.Path[:2].String()
		supported := false
		operationHost := *self.OperationHost
		switch operationHost {
		case "ORCHESTRATOR":
			supported = true
		case "SELF", "HOST":
			if path == "topology_template.node_templates" {
				supported = true
			}
		case "SOURCE", "TARGET":
			if path == "topology_template.relationship_templates" {
				supported = true
			}
		}

		if !supported {
			self.Context.FieldChild("operation_host", operationHost).ReportFieldUnsupportedValue()
		}
	}
}

func (self *InterfaceImplementation) Normalize(o *normal.Operation) {
	if self.Primary != nil {
		o.Implementation = *self.Primary
	}

	if self.Dependencies != nil {
		o.Dependencies = *self.Dependencies
	}

	if self.Timeout != nil {
		o.Timeout = *self.Timeout
	}

	if self.OperationHost != nil {
		o.Host = *self.OperationHost
	}
}
