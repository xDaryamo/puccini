package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// NodeTemplateInstances
//
// [https://cloudify.co/guide/3.2/dsl-spec-node-templates.html]
//

type NodeTemplateInstances struct {
	*Entity `name:"node template instances"`

	Deploy *int64 `read:"deploy"`
}

func NewNodeTemplateInstances(context *tosca.Context) *NodeTemplateInstances {
	return &NodeTemplateInstances{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadNodeTemplateInstances(context *tosca.Context) tosca.EntityPtr {
	self := NewNodeTemplateInstances(context)

	context.ValidateUnsupportedFields(context.ReadFields(self))

	if self.Deploy == nil {
		// Default to 1
		var deploy int64 = 1
		self.Deploy = &deploy
	}

	return self
}
