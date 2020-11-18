package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
)

//
// OperationAssignment
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.15
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.13
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.13
//

// tosca.Reader signature
func ReadOperationAssignment(context *tosca.Context) tosca.EntityPtr {
	context.SetReadTag("Outputs", "")

	return tosca_v2_0.ReadOperationAssignment(context)
}
