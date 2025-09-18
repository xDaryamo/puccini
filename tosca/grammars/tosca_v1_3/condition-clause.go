package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

// Alias: type remains from v2.0, no struct copy needed
type (
	ConditionClause = tosca_v2_0.ConditionClause
)

// ([parsing.Reader] signature) - override of "ConditionClause"
func ReadConditionClause(context *parsing.Context) parsing.EntityPtr {

	// If input is of type "attribute: [ constraints... ]"
	if m, ok := context.Data.(ard.Map); ok && len(m) == 1 {
		for k, v := range m {
			name := yamlkeys.KeyString(k)

			// Only asserts (not and/or/not) have the old "constraints:"
			if name != "and" && name != "or" && name != "not" && name != "assert" {
				if list, ok := v.(ard.List); ok && len(list) > 0 {

					var validation any
					if len(list) == 1 {
						validation = list[0] // single constraint
					} else {
						validation = ard.Map{"$and": list} // multiple constraints → AND
					}
					m[k] = validation // replaces the list with 2.0 constraints
				}
			}
		}
	}

	// Delegate to TOSCA 2.0 logic (which will now find "validation")
	return tosca_v2_0.ReadConditionClause(context).(*tosca_v2_0.ConditionClause)
}

// ([parsing.Reader] signature) - override of "ConditionClauseAnd"
func ReadConditionClauseAnd(context *parsing.Context) parsing.EntityPtr {
	self := tosca_v2_0.NewConditionClause(context)

	if context.ValidateType(ard.TypeList) {
		op := "and"
		self.Operator = &op
		// Use **this** ReadConditionClause (the adapter) for items
		context.ReadListItems(ReadConditionClause, func(item ard.Value) {
			self.ConditionClauses = append(
				self.ConditionClauses,
				item.(*tosca_v2_0.ConditionClause),
			)
		})
	}
	return self
}
