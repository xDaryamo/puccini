package tosca_v2_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

//
// NodeFilter
//
// [TOSCA-v2.0] @ 8.6
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.5
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.5
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.4
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.4
//

type NodeFilter struct {
	*Entity `name:"node filter"`

	// TOSCA 2.0: node_filter is now a single condition clause
	ConditionClause *ValidationClause `read:".,ValidationClause"`

	// Legacy support for TOSCA 1.x (will be nil in TOSCA 2.0)
	PropertyFilters   PropertyFilters   `read:"properties,{}PropertyFilter"`
	CapabilityFilters CapabilityFilters `read:"capabilities,{}CapabilityFilter"`
}

func NewNodeFilter(context *parsing.Context) *NodeFilter {
	return &NodeFilter{
		Entity: NewEntity(context),
	}
}

// ([parsing.Reader] signature)
func ReadNodeFilter(context *parsing.Context) parsing.EntityPtr {
	self := NewNodeFilter(context)

	if context.Is(ard.TypeMap) {
		data := context.Data.(ard.Map)

		// Check if this is TOSCA 2.0 syntax (condition clause) or legacy syntax
		isTosca20ConditionClause := false

		// Look for validation operators (start with '$') at the top level
		for key := range data {
			keyStr := yamlkeys.KeyString(key)
			if len(keyStr) > 0 && keyStr[0] == '$' {
				// This looks like a TOSCA 2.0 condition clause
				isTosca20ConditionClause = true
				break
			}
		}

		// Also check for legacy TOSCA 1.x keys
		hasLegacyKeys := false
		for key := range data {
			keyStr := yamlkeys.KeyString(key)
			if keyStr == "properties" || keyStr == "capabilities" {
				hasLegacyKeys = true
				break
			}
		}

		if isTosca20ConditionClause && !hasLegacyKeys {
			// TOSCA 2.0: Read the entire map as a single condition clause
			self.ConditionClause = ReadValidationClause(context).(*ValidationClause)
		} else if hasLegacyKeys && !isTosca20ConditionClause {
			// Legacy TOSCA 1.x: Read properties and capabilities filters
			context.ValidateUnsupportedFields(context.ReadFields(self))
		} else {
			// Mixed or ambiguous syntax - report error
			context.ReportValueMalformed("node filter", "cannot mix TOSCA 2.0 condition clause syntax with legacy TOSCA 1.x syntax")
		}
	}

	return self
}

func (self *NodeFilter) Normalize(normalRequirement *normal.Requirement) {
	if self.ConditionClause != nil {
		// TOSCA 2.0: Convert condition clause to function call
		fc := self.ConditionClause.ToFunctionCall(self.Context, false)
		NormalizeFunctionCallArguments(fc, self.Context)

		// Store the node filter as a special validation that will be handled by the resolver
		if normalRequirement.NodeTemplatePropertyValidation == nil {
			normalRequirement.NodeTemplatePropertyValidation = make(normal.FunctionCallMap)
		}

		// Use a special key to store the node filter condition
		normalRequirement.NodeTemplatePropertyValidation["$node_filter"] = append(
			normalRequirement.NodeTemplatePropertyValidation["$node_filter"],
			normal.NewFunctionCall(fc),
		)
	} else {
		// Legacy TOSCA 1.x support
		self.PropertyFilters.Normalize(normalRequirement.NodeTemplatePropertyValidation)
		self.CapabilityFilters.Normalize(normalRequirement)
	}
}
