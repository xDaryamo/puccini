package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

// Embeds tosca_v2_0.Schema for compatibility
type Schema struct{ *tosca_v2_0.Schema }

// Allows old code to call this method
func (s *Schema) GetConstraintClauses() ConstraintClauses {
	if s.ValidationClause == nil {
		return nil
	}
	return ConstraintClauses{s.ValidationClause}
}

// Reader for TOSCA 1.3 schema
func ReadSchema(ctx *parsing.Context) parsing.EntityPtr {

	// Convert "constraints" to "validation" for compatibility with 2.0
	if m, ok := ctx.Data.(ard.Map); ok {
		if c, ok := m["constraints"]; ok {
			m["validation"] = c
			delete(m, "constraints")
		}
	}

	// Delegate to the 2.0 reader
	v2s := tosca_v2_0.ReadSchema(ctx).(*tosca_v2_0.Schema)

	// Return the 2.0 entity for compatibility
	return v2s
}
