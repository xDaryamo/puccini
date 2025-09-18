package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// AttributeDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.11
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.10
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.10
//

// ([parsing.Reader] signature)
func ReadAttributeDefinition(context *parsing.Context) parsing.EntityPtr {
	// TOSCA 1.2 doesn't support the "KeySchema" field (introduced in TOSCA 1.3)
	context.SetReadTag("KeySchema", "")
	// TOSCA 1.2 doesn't support the "ValidationClause" field (introduced in TOSCA 2.0)
	context.SetReadTag("ValidationClause", "")

	return tosca_v1_3.ReadAttributeDefinition(context)
}
