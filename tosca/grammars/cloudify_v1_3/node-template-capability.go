package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// NodeTemplateCapability
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-node-templates/]
//

type NodeTemplateCapability struct {
	*Entity `name:"node template capability"`

	Properties Values `read:"properties,Value"`
}

func NewNodeTemplateCapability(context *tosca.Context) *NodeTemplateCapability {
	return &NodeTemplateCapability{
		Entity:     NewEntity(context),
		Properties: make(Values),
	}
}

// tosca.Reader signature
func ReadNodeTemplateCapability(context *tosca.Context) interface{} {
	self := NewNodeTemplateCapability(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

func (self *NodeTemplateCapability) ValidateScalableProperties() {
	for key, value := range self.Properties {
		childContext := self.Context.MapChild(key, value.Data)
		switch key {
		case "default_instances":
			childContext.ValidateType("integer")
		case "min_instances":
			childContext.ValidateType("integer")
		case "max_instances":
			childContext.ValidateType("integer", "string")
		default:
			childContext.ReportFieldUnsupported()
		}
	}
}

//
// NodeTemplateCapabilities
//

type NodeTemplateCapabilities map[string]*NodeTemplateCapability

func (self NodeTemplateCapabilities) Validate() {
	for capabilityName, capability := range self {
		switch capabilityName {
		case "scalable":
			capability.ValidateScalableProperties()
		default:
			capability.Context.ReportFieldUnsupported()
		}
	}
}
