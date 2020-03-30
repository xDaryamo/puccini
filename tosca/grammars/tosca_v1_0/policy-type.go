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
	context.SetReadTag("TriggerDefinitions", "")

	return tosca_v1_3.ReadPolicyType(context)
}
