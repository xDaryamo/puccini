package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

//
// InterfaceDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.16
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.14
//

// tosca.Reader signature
func ReadInterfaceDefinition(context *tosca.Context) interface{} {
	if context.ReadOverrides == nil {
		context.ReadOverrides = make(map[string]string)
	}
	context.ReadOverrides["OperationDefinitions"] = "?,OperationDefinition"
	context.ReadOverrides["NotificationDefinitions"] = ""

	self := tosca_v1_3.NewInterfaceDefinition(context)
	context.ReadFields(self)
	return self
}
