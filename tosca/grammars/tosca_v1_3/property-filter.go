// tosca_v1_3/property_filter_adapter.go
package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

// Inherits from tosca_v2_0.PropertyFilter and provides legacy "Constraint..." support
type PropertyFilter struct {
	*tosca_v2_0.PropertyFilter
}

func (p *PropertyFilter) GetConstraintClauses() ConstraintClauses {
	// Converts ValidationClauses to ConstraintClauses for backward compatibility
	cc := make(ConstraintClauses, len(p.ValidationClauses))
	for i, vc := range p.ValidationClauses {
		cc[i] = (*ConstraintClause)(vc)
	}
	return cc
}

// Reader override for TOSCA 1.3 grammar
func ReadPropertyFilter(ctx *parsing.Context) parsing.EntityPtr {

	// Rename "constraints" to "validation" for compatibility with 2.0
	if m, ok := ctx.Data.(ard.Map); ok {
		if c, ok := m["constraints"]; ok {
			m["validation"] = c
			delete(m, "constraints")
		}
	}

	v2pf := tosca_v2_0.ReadPropertyFilter(ctx).(*tosca_v2_0.PropertyFilter)

	// Return the 1.3 compatible wrapper
	return &PropertyFilter{v2pf}
}
