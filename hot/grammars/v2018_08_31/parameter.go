package v2018_08_31

import (
	"github.com/tliron/puccini/tosca"
)

//
// Parameter
//
// [https://docs.openstack.org/heat/rocky/template_guide/hot_spec.html#parameters-section]
//

type Parameter struct {
	*Entity `name:"parameter"`

	Type        *string `read:"type"`
	Label       *string `read:"label"`
	Description *string `read:"description"`
	Default     *Value  `read:"default,Value"`
	Hidden      *bool   `read:"hidden"`
	Constraints *Value  `read:"constraints,Value"`
	Immutable   *bool   `read:"immutable"`
	Tags        *string `read:"tags"`
}

func NewParameter(context *tosca.Context) *Parameter {
	return &Parameter{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadParameter(context *tosca.Context) interface{} {
	self := NewParameter(context)
	context.ValidateUnsupportedFields(append(context.ReadFields(self, Readers)))
	return self
}
