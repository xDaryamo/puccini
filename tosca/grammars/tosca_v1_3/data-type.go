package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

// Embedding: inherits all methods from the v2.0 struct
type DataType struct {
	*tosca_v2_0.DataType
	// Optional: if old code accesses the field directly
	ConstraintClauses ConstraintClauses `json:"-" yaml:"-"`
}

// ([parsing.Reader] signature)
func ReadDataType(ctx *parsing.Context) parsing.EntityPtr {

	// Transform YAML 1.x: constraints → validation
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

	// Disable TOSCA 2.0 specific fields that require UnitsReader
	ctx.SetReadTag("Units", "")         // disable Units field
	ctx.SetReadTag("CanonicalUnit", "") // disable CanonicalUnit field
	ctx.SetReadTag("Prefixes", "")      // disable Prefixes field
	ctx.SetReadTag("DataTypeName", "")  // disable DataTypeName field

	// Call TOSCA 2.0 reader
	v2dt := tosca_v2_0.ReadDataType(ctx).(*tosca_v2_0.DataType)

	// Return the v2.0 entity directly
	return v2dt
}
