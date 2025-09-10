package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// RelationshipType
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.10
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.10
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.10
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.9
//

// ([parsing.Reader] signature)
func ReadRelationshipType(context *parsing.Context) parsing.EntityPtr {
	// TOSCA 1.3 uses "valid_target_types" while TOSCA 2.0 uses "valid_capability_types"
	// We need to map the v1.3 field name to what the v2.0 reader expects
	context.SetReadTag("ValidCapabilityTypeNames", "valid_target_types")

	return tosca_v2_0.ReadRelationshipType(context)
}
