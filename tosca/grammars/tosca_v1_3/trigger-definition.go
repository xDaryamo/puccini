package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// TriggerDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.22
//

// ([parsing.Reader] signature)
func ReadTriggerDefinition(context *parsing.Context) parsing.EntityPtr {
	// TOSCA 1.3 supports fields that v2.0 doesn't have
	// Remove them temporarily to avoid "unsupported field" errors
	if context.ValidateType(ard.TypeMap) {
		data := context.Data.(ard.Map)

		// Store and remove TOSCA 1.3 specific fields
		var schedule, targetFilter, condition interface{}
		var hasSchedule, hasTargetFilter, hasCondition bool

		if s, ok := data["schedule"]; ok {
			schedule = s
			hasSchedule = true
			delete(data, "schedule")
		}
		if tf, ok := data["target_filter"]; ok {
			targetFilter = tf
			hasTargetFilter = true
			delete(data, "target_filter")
		}
		if c, ok := data["condition"]; ok {
			condition = c
			hasCondition = true
			delete(data, "condition")
		}

		// Call v2.0 reader
		result := tosca_v2_0.ReadTriggerDefinition(context)

		// Restore fields if needed (for completeness)
		if hasSchedule {
			data["schedule"] = schedule
		}
		if hasTargetFilter {
			data["target_filter"] = targetFilter
		}
		if hasCondition {
			data["condition"] = condition
		}

		return result
	}

	return tosca_v2_0.ReadTriggerDefinition(context)
}
