package v1_1

import (
	"github.com/tliron/puccini/tosca"
)

//
// PolicyType
//

type PolicyType struct {
	*Type `name:"policy type"`

	PropertyDefinitions            PropertyDefinitions `read:"properties,PropertyDefinition" inherit:"properties,Parent"`
	TargetNodeTypeOrGroupTypeNames *[]string           `read:"targets" inherit:"targets,Parent"`
	TriggerDefinitions             TriggerDefinitions  `read:"triggers" inherit:"triggers,Parent"`

	Parent           *PolicyType  `lookup:"derived_from,ParentName" json:"-" yaml:"-"`
	TargetNodeTypes  []*NodeType  `lookup:"targets,TargetNodeTypeOrGroupTypeNames" inherit:"targets,Parent" json:"-" yaml:"-"`
	TargetGroupTypes []*GroupType `lookup:"targets,TargetNodeTypeOrGroupTypeNames" inherit:"targets,Parent" json:"-" yaml:"-"`
}

func NewPolicyType(context *tosca.Context) *PolicyType {
	return &PolicyType{
		Type:                NewType(context),
		PropertyDefinitions: make(PropertyDefinitions),
	}
}

// tosca.Reader signature
func ReadPolicyType(context *tosca.Context) interface{} {
	self := NewPolicyType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self, Readers))
	return self
}

// tosca.Hierarchical interface
func (self *PolicyType) GetParent() interface{} {
	return self.Parent
}

// tosca.Inherits interface
func (self *PolicyType) Inherit() {
	log.Infof("{inherit} policy type: %s", self.Name)

	if self.Parent == nil {
		return
	}

	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)
}
