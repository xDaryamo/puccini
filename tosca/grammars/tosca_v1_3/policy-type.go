package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// PolicyType
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.12
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.12
//

type PolicyType struct {
	*Type `name:"policy type"`

	PropertyDefinitions            PropertyDefinitions `read:"properties,PropertyDefinition" inherit:"properties,Parent"`
	TargetNodeTypeOrGroupTypeNames *[]string           `read:"targets" inherit:"targets,Parent"`
	TriggerDefinitions             TriggerDefinitions  `read:"triggers,TriggerDefinition" inherit:"triggers,Parent"`

	Parent           *PolicyType  `lookup:"derived_from,ParentName" json:"-" yaml:"-"`
	TargetNodeTypes  []*NodeType  `lookup:"targets,TargetNodeTypeOrGroupTypeNames" inherit:"targets,Parent" json:"-" yaml:"-"`
	TargetGroupTypes []*GroupType `lookup:"targets,TargetNodeTypeOrGroupTypeNames" inherit:"targets,Parent" json:"-" yaml:"-"`
}

func NewPolicyType(context *tosca.Context) *PolicyType {
	return &PolicyType{
		Type:                NewType(context),
		PropertyDefinitions: make(PropertyDefinitions),
		TriggerDefinitions:  make(TriggerDefinitions),
	}
}

// tosca.Reader signature
func ReadPolicyType(context *tosca.Context) interface{} {
	self := NewPolicyType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
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

//
// PolicyTypes
//

type PolicyTypes []*PolicyType
