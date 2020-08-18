package tosca_v1_1

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
)

//
// InterfaceImplementation
//
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.13.2.3
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.13.2.3
//

// tosca.Reader signature
func ReadInterfaceImplementation(context *tosca.Context) tosca.EntityPtr {
	context.SetReadTag("Timeout", "")
	context.SetReadTag("OperationHost", "")

	return tosca_v2_0.ReadInterfaceImplementation(context)
}
