package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca/parsing"
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

	Parent           *PolicyType `lookup:"derived_from,ParentName" traverse:"ignore" json:"-" yaml:"-"`
	TargetNodeTypes  NodeTypes   `lookup:"targets,TargetNodeTypeOrGroupTypeNames" inherit:"targets,Parent" traverse:"ignore" json:"-" yaml:"-"`
	TargetGroupTypes GroupTypes  `lookup:"targets,TargetNodeTypeOrGroupTypeNames" inherit:"targets,Parent" traverse:"ignore" json:"-" yaml:"-"`
}

func NewPolicyType(context *parsing.Context) *PolicyType {
	return &PolicyType{
		Type:                NewType(context),
		PropertyDefinitions: make(PropertyDefinitions),
		TriggerDefinitions:  make(TriggerDefinitions),
	}
}

// ([parsing.Reader] signature)
func ReadPolicyType(context *parsing.Context) parsing.EntityPtr {
	self := NewPolicyType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// ([parsing.Hierarchical] interface)
func (self *PolicyType) GetParent() parsing.EntityPtr {
	return self.Parent
}

// ([parsing.Inherits] interface)
func (self *PolicyType) Inherit() {
	logInherit.Debugf("policy type: %s", self.Name)

	if self.Parent == nil {
		return
	}

	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)

	// TOSCA spec: "existing trigger definitions may not be changed; new trigger definitions may be added"
	if len(self.TriggerDefinitions) == 0 {
		self.TriggerDefinitions = self.Parent.TriggerDefinitions
	} else {
		// Merge parent triggers with new triggers without overriding existing ones
		for name, parentTrigger := range self.Parent.TriggerDefinitions {
			if _, exists := self.TriggerDefinitions[name]; !exists {
				self.TriggerDefinitions[name] = parentTrigger
			}
		}
	}

	// Note: checking TargetNodeTypeOrGroupTypeNames instead of TargetNodeTypes/TargetGroupTypes
	// because the latter will never be nil
	if self.TargetNodeTypeOrGroupTypeNames == nil {
		self.TargetNodeTypeOrGroupTypeNames = self.Parent.TargetNodeTypeOrGroupTypeNames
		self.TargetNodeTypes = self.Parent.TargetNodeTypes
		self.TargetGroupTypes = self.Parent.TargetGroupTypes
	}
	// Cannot handle the "else" case here because the node type hierarchy may not have been created yet
	// This check is performed in the rendering phase below
}

// ([parsing.Renderable] interface)
func (self *PolicyType) Render() {
	// Avoid rendering more than once
	self.renderOnce.Do(self.render)
}

func (self *PolicyType) render() {
	logRender.Debugf("policy type: %s", self.Name)

	// Note: checking TargetNodeTypeOrGroupTypeNames instead of TargetNodeTypes/TargetGroupTypes
	// because the latter will never be nil
	if self.Parent == nil {
		return
	}

	if self.Parent.TargetNodeTypeOrGroupTypeNames == nil {
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
