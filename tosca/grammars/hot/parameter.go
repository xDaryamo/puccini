package hot

import (
	"github.com/tliron/puccini/tosca"
)

var Types = []string{
	"string",
	"number",
	"json",
	"comma_delimited_list",
	"boolean",
}

func IsTypeValid(type_ string) bool {
	for _, t := range Types {
		if t == type_ {
			return true
		}
	}
	return false
}

//
// Parameter
//
// [https://docs.openstack.org/heat/rocky/template_guide/hot_spec.html#parameters-section]
//

type Parameter struct {
	*Entity `name:"parameter"`

	Type        *string       `read:"type" require:"type"`
	Label       *string       `read:"label"`
	Description *string       `read:"description"`
	Default     *Value        `read:"default,Value"`
	Hidden      *bool         `read:"hidden"`
	Constraints []*Constraint `read:"constraints,[]Constraint"`
	Immutable   *bool         `read:"immutable"`
	Tags        *[]string     `read:"tags"`
}

func NewParameter(context *tosca.Context) *Parameter {
	return &Parameter{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadParameter(context *tosca.Context) interface{} {
	self := NewParameter(context)
	context.ValidateUnsupportedFields(append(context.ReadFields(self, Readers)))

	if (self.Type != nil) && !IsTypeValid(*self.Type) {
		context.FieldChild("type", *self.Type).ReportFieldUnsupportedValue()
	}

	return self
}
