package v2018_08_31

import (
	"github.com/tliron/puccini/tosca"
)

//
// Output
//
// [https://docs.openstack.org/heat/rocky/template_guide/hot_spec.html#outputs-section]
//

type Output struct {
	*Entity `name:"output"`

	Description *string `read:"description"`
	Value       *Value  `read:"value,Value"`
	Condition   *string `read:"condition"`
}

func NewOutput(context *tosca.Context) *Output {
	return &Output{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadOutput(context *tosca.Context) interface{} {
	self := NewOutput(context)
	context.ValidateUnsupportedFields(append(context.ReadFields(self, Readers)))
	return self
}
