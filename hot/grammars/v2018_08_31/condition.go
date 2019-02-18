package v2018_08_31

import (
	"github.com/tliron/puccini/tosca"
)

//
// Condition
//
// [https://docs.openstack.org/heat/rocky/template_guide/hot_spec.html#conditions-section]
//

type Condition struct {
	*Entity `name:"condition"`
}

func NewCondition(context *tosca.Context) *Condition {
	return &Condition{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadCondition(context *tosca.Context) interface{} {
	self := NewCondition(context)
	context.ValidateUnsupportedFields(append(context.ReadFields(self, Readers)))
	return self
}
