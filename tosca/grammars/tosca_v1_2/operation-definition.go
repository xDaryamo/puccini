package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// OperationDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.15
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.13
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.13
//

// ([parsing.Reader] signature)
func ReadOperationDefinition(context *parsing.Context) parsing.EntityPtr {
	// TOSCA 1.2 doesn't support the "OutputDefinitions" field that exists in TOSCA 2.0
	// This disables reading of the OutputDefinitions field
	context.SetReadTag("OutputDefinitions", "")

	return tosca_v2_0.ReadOperationDefinition(context)
}
