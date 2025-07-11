// tosca_v1_3/property_definition_adapter.go
package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

// Embeds all fields/methods from tosca_v2_0.PropertyDefinition
type PropertyDefinition struct {
	*tosca_v2_0.PropertyDefinition
}

// Bridge method: allows 1.3 code to call this
func (p *PropertyDefinition) GetConstraintClauses() ConstraintClauses {
	if p.ValidationClause == nil {
		return nil
	}
	return ConstraintClauses{p.ValidationClause}
}

// [parsing.Reader] signature
func ReadPropertyDefinition(ctx *parsing.Context) parsing.EntityPtr {

	// Convert "constraints" list (1.x) to "validation" (2.0)
	if m, ok := ctx.Data.(ard.Map); ok {
		if c, ok := m["constraints"].(ard.List); ok && len(c) > 0 {
			if len(c) == 1 {
				m["validation"] = c[0]
			} else {
				m["validation"] = ard.Map{"$and": c}
			}
			delete(m, "constraints")
		}
	}

	// Metadata supported in TOSCA 1.3
	// ctx.SetReadTag("Metadata", "") // Removed: metadata is supported in 1.3

	// Use the tosca_v2_0 parser
	v2prop := tosca_v2_0.ReadPropertyDefinition(ctx).(*tosca_v2_0.PropertyDefinition)

	// Return the v2.0 entity to match NodeType.Properties type
	return v2prop
}
