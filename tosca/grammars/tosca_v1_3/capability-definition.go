package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// CapabilityDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.2
//

// ([parsing.Reader] signature) - override
func ReadCapabilityDefinition(context *parsing.Context) parsing.EntityPtr {
	self := tosca_v2_0.NewCapabilityDefinition(context)

	if context.Is(ard.TypeMap) {
		// Transform TOSCA 1.3 specific fields before reading
		if m, ok := context.Data.(ard.Map); ok {
			// Map valid_source_types to valid_source_node_types for compatibility with TOSCA 2.0
			if validSourceTypes, exists := m["valid_source_types"]; exists {
				m["valid_source_node_types"] = validSourceTypes
				delete(m, "valid_source_types")
			}

			// Handle occurrences field - TOSCA 1.3 specific, ignore it for now
			if _, exists := m["occurrences"]; exists {
				delete(m, "occurrences")
			}
		}

		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.CapabilityTypeName = context.FieldChild("type", context.Data).ReadString()
	}

	return self
}
