package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
)

//
// WorkflowActivityCallOperation
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.19.2.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.17.2.3
//

// tosca.Reader signature
func ReadWorkflowActivityCallOperation(context *tosca.Context) tosca.EntityPtr {
	self := tosca_v2_0.NewWorkflowActivityCallOperation(context)
	self.InterfaceAndOperation = context.FieldChild("operation", context.Data).ReadString()
	return self
}
