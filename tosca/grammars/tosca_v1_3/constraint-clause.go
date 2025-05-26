package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

// Instead of creating remapped scriptlets, use the original validation scriptlets directly
// This ensures the keys match what tosca_v2_0.ReadValidationClause expects

var ConstraintClauseScriptlets = tosca_v2_0.ValidationClauseScriptlets
var ConstraintClauseNativeArgumentIndexes = tosca_v2_0.ValidationClauseNativeArgumentIndexes

// Override the reader to ensure it has access to the correct operators
func ReadConstraintClause(context *parsing.Context) parsing.EntityPtr {
	// Create a new context with the right scriptlet namespace
	contextWithScriptlets := context.Clone(context.Data)
	contextWithScriptlets.ScriptletNamespace = DefaultScriptletNamespace

	// Call the v2.0 implementation with our enhanced context
	return tosca_v2_0.ReadValidationClause(contextWithScriptlets)
}
