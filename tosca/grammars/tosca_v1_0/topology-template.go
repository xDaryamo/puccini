package tosca_v1_0

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

//
// TopologyTemplate
//
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.8
//

// tosca.Reader signature
func ReadTopologyTemplate(context *tosca.Context) interface{} {
	if context.ReadOverrides == nil {
		context.ReadOverrides = make(map[string]string)
	}
	context.ReadOverrides["WorkflowDefinitions"] = ""

	return tosca_v1_3.ReadTopologyTemplate(context)
}
