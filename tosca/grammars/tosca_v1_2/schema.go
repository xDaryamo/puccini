package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Schema
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ ?
// [TOSCA-Simple-Profile-YAML-v1.1] @ ?
// [TOSCA-Simple-Profile-YAML-v1.0] @ ?
//

// ([parsing.Reader] signature)
func ReadSchema(context *parsing.Context) parsing.EntityPtr {
	// TOSCA 1.2 doesn't support the "ValidationClause" field (introduced in TOSCA 2.0)
	context.SetReadTag("ValidationClause", "")
	// TOSCA 1.2 doesn't support the "KeySchema" field (introduced in TOSCA 1.3)
	context.SetReadTag("KeySchema", "")
	// TOSCA 1.2 doesn't support the "Metadata" field (introduced in TOSCA 1.3)
	context.SetReadTag("Metadata", "")

	return tosca_v2_0.ReadSchema(context)
}