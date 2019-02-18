package v2018_08_31

import (
	"github.com/tliron/puccini/tosca"
)

//
// Resource
//
// [https://docs.openstack.org/heat/rocky/template_guide/hot_spec.html#resources-section]
//

type Resource struct {
	*Entity `name:"resource"`

	Type           *string `read:"type"`
	Properties     Values  `read:"properties,Value"`
	Metadata       *string `read:"metadata"`
	DependsOn      *string `read:"depends_on"`
	UpdatePolicy   *string `read:"update_policy"`
	DeletionPolicy *string `read:"deletion_policy"`
	ExternalID     *string `read:"external_id"`
	Condition      *string `read:"condition"`
}

func NewResource(context *tosca.Context) *Resource {
	return &Resource{
		Entity:     NewEntity(context),
		Properties: make(Values),
	}
}

// tosca.Reader signature
func ReadResource(context *tosca.Context) interface{} {
	self := NewResource(context)
	context.ValidateUnsupportedFields(append(context.ReadFields(self, Readers)))
	return self
}
