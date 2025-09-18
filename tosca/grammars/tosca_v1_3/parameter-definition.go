package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

// ParameterDefinition wraps the v2.0 implementation but converts constraints to validation
type ParameterDefinition struct {
	*tosca_v2_0.ParameterDefinition
}

func NewParameterDefinition(context *parsing.Context) *ParameterDefinition {
	return &ParameterDefinition{
		ParameterDefinition: tosca_v2_0.NewParameterDefinition(context),
	}
}

// ([parsing.Reader] signature)
func ReadParameterDefinition(context *parsing.Context) parsing.EntityPtr {
	// Convert "constraints" list (1.x) to "validation" (2.0) before calling v2.0 reader
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
							if k == "valid_values" && i == 0 && len(c) == 1 {
								// For single valid_values constraint, ensure array structure is preserved
								// Force the function to use legacy syntax (2 arguments) instead of TOSCA 2.0 syntax (3 arguments)
								if validValuesArray, ok := v.(ard.List); ok {
									// Create a validation clause that will be called with 2 arguments
									// valid_values(currentValue, validValuesArray)
									newConstraintMap[k] = validValuesArray
								} else {
									newConstraintMap[k] = v
								}
							} else {
								newConstraintMap[k] = v
							}
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

	// Use the tosca_v2_0 reader directly and return the v2.0 type
	return tosca_v2_0.ReadParameterDefinition(context)
}

func (self *ParameterDefinition) Inherit(parentDefinition *ParameterDefinition) {
	self.ParameterDefinition.Inherit(parentDefinition.ParameterDefinition)
}
