package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
)

//
// PolicyType
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.12
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.12
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.12
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.11
//

type PolicyType struct {
	*Type `name:"policy type"`

	PropertyDefinitions            PropertyDefinitions `read:"properties,PropertyDefinition" inherit:"properties,Parent"`
	TargetNodeTypeOrGroupTypeNames *[]string           `read:"targets" inherit:"targets,Parent"`
	TriggerDefinitions             TriggerDefinitions  `read:"triggers,TriggerDefinition" inherit:"triggers,Parent"` // introduced in TOSCA 1.1

	Parent           *PolicyType `lookup:"derived_from,ParentName" json:"-" yaml:"-"`
	TargetNodeTypes  NodeTypes   `lookup:"targets,TargetNodeTypeOrGroupTypeNames" inherit:"targets,Parent" json:"-" yaml:"-"`
	TargetGroupTypes GroupTypes  `lookup:"targets,TargetNodeTypeOrGroupTypeNames" inherit:"targets,Parent" json:"-" yaml:"-"`
}

func NewPolicyType(context *tosca.Context) *PolicyType {
	return &PolicyType{
		Type:                NewType(context),
		PropertyDefinitions: make(PropertyDefinitions),
		TriggerDefinitions:  make(TriggerDefinitions),
	}
}

// tosca.Reader signature
func ReadPolicyType(context *tosca.Context) tosca.EntityPtr {
	self := NewPolicyType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Hierarchical interface
func (self *PolicyType) GetParent() tosca.EntityPtr {
	return self.Parent
}

// tosca.Inherits interface
func (self *PolicyType) Inherit() {
	logInherit.Debugf("policy type: %s", self.Name)

	if self.Parent == nil {
		return
	}

	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)

	// TODO: inherit TriggerDefinitions?

	// (Note we are checking for TargetNodeTypeOrGroupTypeNames and not TargetNodeTypes/TargetGroupTypes, because the latter will never be nil)
	if self.TargetNodeTypeOrGroupTypeNames == nil {
		self.TargetNodeTypeOrGroupTypeNames = self.Parent.TargetNodeTypeOrGroupTypeNames
		self.TargetNodeTypes = self.Parent.TargetNodeTypes
		self.TargetGroupTypes = self.Parent.TargetGroupTypes
	}
	// We cannot handle the "else" case here, because the node type hierarchy may not have been created yet,
	// So we will do that check in the rendering phase, below
}

// parser.Renderable interface
func (self *PolicyType) Render() {
	logRender.Debugf("policy type: %s", self.Name)

	// (Note we are checking for TargetNodeTypeOrGroupTypeNames and not TargetNodeTypes/TargetGroupTypes, because the latter will never be nil)
	if (self.Parent == nil) || (self.Parent.TargetNodeTypeOrGroupTypeNames == nil) {
		return
	}

	context := self.Context.FieldChild("targets", nil)
	self.Parent.TargetNodeTypes.ValidateSubset(self.TargetNodeTypes, context)
	self.Parent.TargetGroupTypes.ValidateSubset(self.TargetGroupTypes, context)
}

//
// PolicyTypes
//

type PolicyTypes []*PolicyType
