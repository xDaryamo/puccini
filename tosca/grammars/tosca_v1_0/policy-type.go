package tosca_v1_0

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

//
// PolicyType
//
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.11
//

// tosca.Reader signature
func ReadPolicyType(context *tosca.Context) interface{} {
	if context.ReadOverrides == nil {
		context.ReadOverrides = make(map[string]string)
	}
	context.ReadOverrides["TriggerDefinitions"] = ""

	return tosca_v1_3.ReadPolicyType(context)
}
