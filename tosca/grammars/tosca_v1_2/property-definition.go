package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

//
// PropertyDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.9
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.8
//

// tosca.Reader signature
func ReadPropertyDefinition(context *tosca.Context) interface{} {
	if context.ReadOverrides == nil {
		context.ReadOverrides = make(map[string]string)
	}
	context.ReadOverrides["Metadata"] = ""
	context.ReadOverrides["KeySchema"] = ""

	return tosca_v1_3.ReadPropertyDefinition(context)
}
