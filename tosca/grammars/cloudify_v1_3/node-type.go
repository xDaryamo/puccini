package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// NodeType
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-node-types/]
//

type NodeType struct {
	*Type `name:"node type"`

	InterfaceDefinitions InterfaceDefinitions `read:"interfaces,InterfaceDefinition" inherit:"interfaces,Parent"`
	PropertyDefinitions  PropertyDefinitions  `read:"properties,PropertyDefinition" inherit:"properties,Parent"`

	Parent *NodeType `lookup:"derived_from,ParentName" json:"-" yaml:"-"`
}

func NewNodeType(context *tosca.Context) *NodeType {
	return &NodeType{
		Type:                 NewType(context),
		InterfaceDefinitions: make(InterfaceDefinitions),
		PropertyDefinitions:  make(PropertyDefinitions),
	}
}

// tosca.Reader signature
func ReadNodeType(context *tosca.Context) tosca.EntityPtr {
	self := NewNodeType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Hierarchical interface
func (self *NodeType) GetParent() tosca.EntityPtr {
	return self.Parent
}

// tosca.Inherits interface
func (self *NodeType) Inherit() {
	logInherit.Debugf("node type: %s", self.Name)

	if self.Parent == nil {
		return
	}

	self.InterfaceDefinitions.Inherit(self.Parent.InterfaceDefinitions)
	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)
}

//
// NodeTypes
//

type NodeTypes []*NodeType
