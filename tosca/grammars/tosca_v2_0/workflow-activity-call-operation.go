package tosca_v2_0

import (
	"strings"

	"github.com/tliron/puccini/tosca"
)

//
// WorkflowActivityCallOperation
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.23.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.19.2.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.17.2.3
//

type WorkflowActivityCallOperation struct {
	*Entity `name:"workflow activity call operation"`
	Name    string

	CallOperationSpec *string

	CallInterface *InterfaceAssignment `json:"-" yaml:"-"`
	CallOperation *OperationAssignment `json:"-" yaml:"-"`
}

func NewWorkflowActivityCallOperation(context *tosca.Context) *WorkflowActivityCallOperation {
	return &WorkflowActivityCallOperation{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadWorkflowActivityCallOperation(context *tosca.Context) tosca.EntityPtr {
	self := NewWorkflowActivityCallOperation(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

func (self *WorkflowActivityCallOperation) Render(stepDefinition *WorkflowStepDefinition) {
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
