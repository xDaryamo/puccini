package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

//
// CapabilityAssignment
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.1
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.1
//

// tosca.Reader signature
func ReadCapabilityAssignment(context *tosca.Context) interface{} {
	context.SetReadTag("Occurrences", "")

	return tosca_v1_3.ReadCapabilityAssignment(context)
}
