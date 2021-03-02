package tosca_v2_0

import (
	"strings"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/yamlkeys"
)

//
// WorkflowActivityDefinition
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.23
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.19
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.17
//

type WorkflowActivityDefinition struct {
	*Entity `name:"workflow activity definition"`

	DelegateWorkflowDefinitionName *string
	InlineWorkflowDefinitionName   *string
	SetNodeState                   *string
	CallOperationSpec              *string

	DelegateWorkflowDefinition *WorkflowDefinition  `lookup:"delegate,DelegateWorkflowDefinitionName" json:"-" yaml:"-"`
	InlineWorkflowDefinition   *WorkflowDefinition  `lookup:"inline,InlineWorkflowDefinitionName" json:"-" yaml:"-"`
	CallInterface              *InterfaceAssignment `json:"-" yaml:"-"`
	CallOperation              *OperationAssignment `json:"-" yaml:"-"`
}

func NewWorkflowActivityDefinition(context *tosca.Context) *WorkflowActivityDefinition {
	return &WorkflowActivityDefinition{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadWorkflowActivityDefinition(context *tosca.Context) tosca.EntityPtr {
	self := NewWorkflowActivityDefinition(context)

	if context.ValidateType(ard.TypeMap) {
		map_ := context.Data.(ard.Map)
		if len(map_) != 1 {
			context.ReportValueMalformed("workflow activity definition", "map length not 1")
			return self
		}

		for key, value := range map_ {
			operator := yamlkeys.KeyString(key)
			childContext := context.FieldChild(operator, value)

			switch operator {
			case "delegate":
				self.DelegateWorkflowDefinitionName = childContext.ReadString()
			case "inline":
				self.InlineWorkflowDefinitionName = childContext.ReadString()
			case "set_state":
				self.SetNodeState = childContext.ReadString()
			case "call_operation":
				self.CallOperationSpec = childContext.ReadString()
			default:
				context.ReportValueMalformed("workflow activity definition", "unsupported operator")
				return self
			}

			// We have only one key
			break
		}
	}

	return self
}

func (self *WorkflowActivityDefinition) Render(stepDefinition *WorkflowStepDefinition) {
	if self.CallOperationSpec == nil {
		return
	}

	// Parse operation spec
	s := strings.SplitN(*self.CallOperationSpec, ".", 2)
	if len(s) != 2 {
		self.Context.FieldChild("call_operation", *self.CallOperationSpec).ReportValueWrongFormat("interface.operation")
		return
	}

	var ok bool

	// Lookup interface by name
	if stepDefinition.TargetNodeTemplate != nil {
		if self.CallInterface, ok = stepDefinition.TargetNodeTemplate.Interfaces[s[0]]; !ok {
			self.Context.FieldChild("call_operation", s[0]).ReportReferenceNotFound("interface", stepDefinition.TargetNodeTemplate)
			return
		}
	} else if stepDefinition.TargetGroup != nil {
		if self.CallInterface, ok = stepDefinition.TargetGroup.Interfaces[s[0]]; !ok {
			self.Context.FieldChild("call_operation", s[0]).ReportReferenceNotFound("interface", stepDefinition.TargetGroup)
			return
		}
	} else {
		// There was a lookup problem (neither node template nor group)
		return
	}

	// Lookup operation by name
	if self.CallOperation, ok = self.CallInterface.Operations[s[1]]; !ok {
		self.Context.FieldChild("call_operation", s[1]).ReportReferenceNotFound("operation", self.CallInterface)
	}
}

func (self *WorkflowActivityDefinition) Normalize(normalWorkflowStep *normal.WorkflowStep) *normal.WorkflowActivity {
	logNormalize.Debug("workflow activity")

	normalWorkflowActivity := normalWorkflowStep.NewActivity()
	if self.DelegateWorkflowDefinition != nil {
		normalWorkflowActivity.DelegateWorkflow = normalWorkflowStep.Workflow.ServiceTemplate.Workflows[self.DelegateWorkflowDefinition.Name]
	} else if self.InlineWorkflowDefinition != nil {
		normalWorkflowActivity.InlineWorkflow = normalWorkflowStep.Workflow.ServiceTemplate.Workflows[self.InlineWorkflowDefinition.Name]
	} else if self.SetNodeState != nil {
		normalWorkflowActivity.SetNodeState = *self.SetNodeState
	} else if self.CallOperation != nil {
		var normalInterface *normal.Interface
		if normalWorkflowStep.TargetNodeTemplate != nil {
			normalInterface = normalWorkflowStep.TargetNodeTemplate.Interfaces[self.CallInterface.Name]
		} else if normalWorkflowStep.TargetGroup != nil {
			normalInterface = normalWorkflowStep.TargetGroup.Interfaces[self.CallInterface.Name]
		} else {
			return normalWorkflowActivity
		}
		normalWorkflowActivity.CallOperation = normalInterface.Operations[self.CallOperation.Name]
	}

	return normalWorkflowActivity
}

//
// WorkflowActivityDefinitions
//

type WorkflowActivityDefinitions []*WorkflowActivityDefinition
