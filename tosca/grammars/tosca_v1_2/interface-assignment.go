package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

//
// InterfaceAssignment
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.16
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.14
//

// tosca.Reader signature
func ReadInterfaceAssignment(context *tosca.Context) interface{} {
	if context.ReadOverrides == nil {
		context.ReadOverrides = make(map[string]string)
	}
	context.ReadOverrides["Operations"] = "?,OperationAssignment"
	context.ReadOverrides["Notifications"] = ""

	self := tosca_v1_3.NewInterfaceAssignment(context)
	context.ReadFields(self)
	return self
}
