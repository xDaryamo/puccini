package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

//
// InterfaceType
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.5
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.5
//

// tosca.Reader signature
func ReadInterfaceType(context *tosca.Context) interface{} {
	if context.ReadOverrides == nil {
		context.ReadOverrides = make(map[string]string)
	}
	context.ReadOverrides["OperationDefinitions"] = "?,OperationDefinition"
	context.ReadOverrides["NotificationDefinitions"] = ""

	self := tosca_v1_3.NewInterfaceType(context)
	context.ReadFields(self)
	return self
}
