package tosca_v1_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// RelationshipType
//
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.6
//

// ([parsing.Reader] signature)
func ReadRelationshipType(context *parsing.Context) parsing.EntityPtr {
	// Convert "valid_target_types" (1.0) to "valid_target_node_types" (2.0) before calling v2.0 reader
	if context.Is(ard.TypeMap) {
		if m, ok := context.Data.(ard.Map); ok {
			if validTargetTypes, ok := m["valid_target_types"]; ok {
				// In TOSCA 1.0, valid_target_types refers to capability types
				// In TOSCA 2.0, this maps to valid_capability_types
				m["valid_capability_types"] = validTargetTypes
				delete(m, "valid_target_types")
			}
		}
	}

	return tosca_v2_0.ReadRelationshipType(context)
}