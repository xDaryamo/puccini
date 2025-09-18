package tosca_v1_1

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// PropertyDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.8
//

// ([parsing.Reader] signature)
func ReadPropertyDefinition(context *parsing.Context) parsing.EntityPtr {
	// Convert "constraints" list (1.1) to "validation" (2.0) before calling v2.0 reader
	if context.Is(ard.TypeMap) {
		if m, ok := context.Data.(ard.Map); ok {
			if c, ok := m["constraints"].(ard.List); ok && len(c) > 0 {
				// Convert constraints array to validation clause
				processedConstraints := make(ard.List, len(c))
				for i, constraint := range c {
					if constraintMap, ok := constraint.(ard.Map); ok {
						// Make a copy to avoid modifying the original
						newConstraintMap := make(ard.Map)
						for k, v := range constraintMap {
							newConstraintMap[k] = v
						}
						processedConstraints[i] = newConstraintMap
					} else {
						processedConstraints[i] = constraint
					}
				}

				if len(processedConstraints) == 1 {
					m["validation"] = processedConstraints[0]
				} else {
					m["validation"] = ard.Map{"$and": processedConstraints}
				}
				delete(m, "constraints")
			}
		}
	}

	// TOSCA 1.1 doesn't support the "KeySchema" field (introduced in TOSCA 1.3)
	context.SetReadTag("KeySchema", "")

	return tosca_v2_0.ReadPropertyDefinition(context)
}