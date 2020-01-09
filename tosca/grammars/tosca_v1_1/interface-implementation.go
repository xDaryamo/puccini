package tosca_v1_1

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

//
// InterfaceImplementation
//
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.13.2.3
//

// tosca.Reader signature
func ReadInterfaceImplementation(context *tosca.Context) interface{} {
	if context.ReadOverrides == nil {
		context.ReadOverrides = make(map[string]string)
	}
	context.ReadOverrides["Timeout"] = ""
	context.ReadOverrides["OperationHost"] = ""

	return tosca_v1_3.ReadInterfaceImplementation(context)
}
