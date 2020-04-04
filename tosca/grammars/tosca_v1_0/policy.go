package tosca_v1_0

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

//
// Policy
//
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.6
//

// tosca.Reader signature
func ReadPolicy(context *tosca.Context) tosca.EntityPtr {
	context.SetReadTag("TriggerDefinitions", "")

	return tosca_v1_3.ReadPolicy(context)
}
