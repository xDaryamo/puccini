package tosca_v2_0

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

// =============================================================================
// MAIN STRUCT DEFINITION
// =============================================================================
//
// SubstitutionMappings represents the TOSCA substitution mappings grammar
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.13, 2.10, 2.11, 2.12
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.12, 2.10, 2.11
// [TOSCA-Simple-Profile-YAML-v1.1] @ 2.10, 2.11
// [TOSCA-Simple-Profile-YAML-v1.0] @ 2.10, 2.11
//

type SubstitutionMappings struct {
	*Entity `name:"substitution mappings"`

	NodeTypeName        *string             `read:"node_type" mandatory:""`
	CapabilityMappings  CapabilityMappings  `read:"capabilities,CapabilityMapping"`
	RequirementMappings RequirementMappings `read:"requirements,RequirementMapping"`
	PropertyMappings    PropertyMappings    `read:"properties,PropertyMapping"`     // introduced in TOSCA 1.2
	AttributeMappings   AttributeMappings   `read:"attributes,AttributeMapping"`    // introduced in TOSCA 1.3
	InterfaceMappings   InterfaceMappings   `read:"interfaces,InterfaceMapping"`    // introduced in TOSCA 1.2
	SubstitutionFilter  *NodeFilter         `read:"substitution_filter,NodeFilter"` // introduced in TOSCA 1.3

	NodeType *NodeType `lookup:"node_type,NodeTypeName" traverse:"ignore" json:"-" yaml:"-"`
}

// =============================================================================
// CONSTRUCTOR FUNCTIONS
// =============================================================================

func NewSubstitutionMappings(context *parsing.Context) *SubstitutionMappings {
	return &SubstitutionMappings{
		Entity:              NewEntity(context),
		CapabilityMappings:  make(CapabilityMappings),
		RequirementMappings: make(RequirementMappings),
		PropertyMappings:    make(PropertyMappings),
		AttributeMappings:   make(AttributeMappings),
		InterfaceMappings:   make(InterfaceMappings),
	}
}

// =============================================================================
// PARSING CONSTANTS AND HELPER FUNCTIONS
// =============================================================================

const (
	requirementMappingStr = "requirement mapping"
	selectDirective       = "select"
)

// =============================================================================
// PARSING HELPER FUNCTIONS - TYPE DETECTION
// =============================================================================

// isListOfListsForRequirement checks if a list contains only lists of exactly 2 strings
func isListOfListsForRequirement(list ard.List) bool {
	if len(list) == 0 {
		return false
	}

	// Check if all elements are lists of exactly 2 strings
	for _, item := range list {
		if itemList, ok := item.(ard.List); ok {
			if len(itemList) != 2 {
				return false
			}
			// Check if both elements are strings
			if _, ok1 := itemList[0].(string); !ok1 {
				return false
			}
			if _, ok2 := itemList[1].(string); !ok2 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

// isListOfStringPairs returns true if the list is a list of [string, string] pairs
func isListOfStringPairs(list ard.List) bool {
	for _, item := range list {
		pair, ok := item.(ard.List)
		if !ok || len(pair) != 2 {
			return false
		}
		if _, ok := pair[0].(string); !ok {
			return false
		}
		if _, ok := pair[1].(string); !ok {
			return false
		}
	}
	return true
}

// isListOfStrings returns true if the list is a list of strings (selectable nodes)
func isListOfStrings(list ard.List) bool {
	for _, item := range list {
		if _, ok := item.(string); !ok {
			return false
		}
	}
	return true
}

// isMixedFormatList checks if a list contains both traditional mappings and selectable nodes
func isMixedFormatList(list ard.List) bool {
	hasLists := false
	hasStrings := false

	for _, item := range list {
		if itemList, ok := item.(ard.List); ok {
			// Check if it's a valid 2-element list
			if len(itemList) == 2 {
				if _, ok1 := itemList[0].(string); ok1 {
					if _, ok2 := itemList[1].(string); ok2 {
						hasLists = true
					}
				}
			}
		} else if _, ok := item.(string); ok {
			hasStrings = true
		}
	}

	// Mixed format if we have both lists and strings
	return hasLists && hasStrings
}

// =============================================================================
// PARSING HELPER FUNCTIONS - KEY PROCESSING
// =============================================================================

// isSerializedListKey checks if a string is a serialized list key like "[ service, 2 ]"
func (self *SubstitutionMappings) isSerializedListKey(keyStr string) bool {
	trimmed := strings.TrimSpace(keyStr)
	return strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]")
}

// parseSerializedListKey parses a serialized list key like "[ service, 2 ]" or "[ service, UNBOUNDED ]"
func (self *SubstitutionMappings) parseSerializedListKey(keyStr string) (reqName string, count string, ok bool) {
	trimmed := strings.TrimSpace(keyStr)
	if !strings.HasPrefix(trimmed, "[") || !strings.HasSuffix(trimmed, "]") {
		return "", "", false
	}

	// Remove brackets
	inner := strings.TrimSpace(trimmed[1 : len(trimmed)-1])

	// Split by comma
	parts := strings.Split(inner, ",")
	if len(parts) != 2 {
		return "", "", false
	}

	reqName = strings.TrimSpace(parts[0])
	count = strings.TrimSpace(parts[1])

	// Remove quotes if present
	if len(reqName) >= 2 && reqName[0] == '"' && reqName[len(reqName)-1] == '"' {
		reqName = reqName[1 : len(reqName)-1]
	}
	if len(count) >= 2 && count[0] == '"' && count[len(count)-1] == '"' {
		count = count[1 : len(count)-1]
	}

	return reqName, count, true
}

// extractKeyString extracts the string representation from various key types including YAMLKey
func (self *SubstitutionMappings) extractKeyString(key interface{}) string {
	switch k := key.(type) {
	case string:
		return k
	default:
		// Handle YAMLKey or other types by converting to string
		keyStr := fmt.Sprintf("%v", k)

		// YAMLKey often formats as "- service\n- 2" or similar
		// Convert this to the bracket format we expect: "[ service, 2 ]"
		if strings.Contains(keyStr, "\n") {
			lines := strings.Split(keyStr, "\n")
			if len(lines) == 2 {
				// Extract the values (remove "- " prefix)
				line1 := strings.TrimSpace(lines[0])
				line2 := strings.TrimSpace(lines[1])

				if strings.HasPrefix(line1, "- ") && strings.HasPrefix(line2, "- ") {
					val1 := strings.TrimPrefix(line1, "- ")
					val2 := strings.TrimPrefix(line2, "- ")

					// Convert to bracket format
					return fmt.Sprintf("[ %s, %s ]", val1, val2)
				}
			}
		}

		// If it's already in bracket format or other format, return as is
		return keyStr
	}
}

// =============================================================================
// MAIN PARSING LOGIC
// =============================================================================

// parseRequirementMapping is the main entry point for parsing requirement mappings
func parseRequirementMapping(self *SubstitutionMappings, context *parsing.Context, reqName string, reqValue interface{}) {
	// If the value is a string, it's a single selectable node (requires select directive)
	switch v := reqValue.(type) {
	case string:
		handleSelectableNodesList(self, context, reqName, ard.List{v})
	case ard.List:
		if len(v) == 2 {
			// If both elements are strings, treat as [node_template, requirement] mapping (never selectable node)
			if _, ok0 := v[0].(string); ok0 {
				if _, ok1 := v[1].(string); ok1 {
					handleTraditionalRequirementMapping(self, context, reqName, v)
					return
				}
			}
		}
		if isListOfStrings(v) {
			handleSelectableNodesList(self, context, reqName, v)
			return
		}
		if isListOfStringPairs(v) {
			handleMultipleRequirementAssignments(self, context, reqName, v)
			return
		}
		context.ReportValueMalformed("requirement mapping", "invalid list format for requirement mapping")
	default:
		context.ReportValueMalformed("requirement mapping", "unsupported type for requirement mapping")
	}
}

// ([parsing.Reader] signature)
func ReadSubstitutionMappings(context *parsing.Context) parsing.EntityPtr {
	if context.HasQuirk(parsing.QuirkSubstitutionMappingsRequirementsList) {
		if map_, ok := context.Data.(ard.Map); ok {
			if requirements, ok := map_["requirements"]; ok {
				if _, ok := requirements.(ard.List); ok {
					context.SetReadTag("RequirementMappings", "requirements,{}RequirementMapping")
				}
			}
		}
	}

	self := NewSubstitutionMappings(context)

	// Handle special requirement formats (TOSCA 2.0) BEFORE normal processing
	if map_, ok := context.Data.(ard.Map); ok {
		if requirements, ok := map_["requirements"]; ok {
			if reqMap, ok := requirements.(ard.Map); ok {
				self.processRequirementMappings(context, reqMap)
			} else if reqList, ok := requirements.(ard.List); ok {
				self.processRequirementMappingsList(context, reqList)
			}
		}
	}

	context.ValidateUnsupportedFields(context.ReadFields(self))

	// Post-process any compact form with count mappings
	self.postProcessCompactFormWithCount(context)

	return self
}

// =============================================================================
// REQUIREMENT MAPPING PROCESSING FUNCTIONS
// =============================================================================

// processRequirementMappings handles requirement mappings from a map structure
func (self *SubstitutionMappings) processRequirementMappings(context *parsing.Context, reqMap ard.Map) {
	processedRequirements := make(map[string]bool)

	for reqName, reqValue := range reqMap {
		reqNameStr := reqName.(string)

		// Check for compact form with count FIRST (highest priority)
		if reqList, ok := reqValue.(ard.List); ok {
			if self.handleCompactFormWithCount(context, reqNameStr, reqList) {
				processedRequirements[reqNameStr] = true
				continue
			}
		}

		// Then check for other list formats
		if reqList, ok := reqValue.(ard.List); ok {
			if self.processListFormats(context, reqNameStr, reqList) {
				processedRequirements[reqNameStr] = true
			}
		}
	}

	// Remove processed requirements from the map to prevent double processing
	if len(processedRequirements) > 0 {
		self.cleanupProcessedRequirements(context, reqMap, processedRequirements)
	}
}

// processRequirementMappingsList handles requirement mappings from a list structure
func (self *SubstitutionMappings) processRequirementMappingsList(context *parsing.Context, reqList ard.List) {
	processedIndices := make(map[int]bool)
	// Track counters for each requirement name to ensure sequential indexing
	requirementCounters := make(map[string]int)

	for i, reqItem := range reqList {
		if processedIndices[i] {
			continue
		}

		if reqMap, ok := reqItem.(ard.Map); ok {
			// Process each key-value pair in the map
			for reqKey, reqValue := range reqMap {
				self.processMixedRequirementMappingWithCounter(context, reqKey, reqValue, i, requirementCounters)
				processedIndices[i] = true
			}
		} else {
			// Handle non-map items (though shouldn't happen in valid TOSCA)
			context.ReportValueMalformed("requirement mapping",
				fmt.Sprintf("requirement item %d is not a map", i))
		}
	}

	// Clear the original requirements list since we've processed everything
	self.clearOriginalRequirements(context)
}

// processListFormats processes different list formats for requirement mappings
func (self *SubstitutionMappings) processListFormats(context *parsing.Context, reqName string, reqList ard.List) bool {
	// Check for mixed format FIRST (contains both traditional mappings and selectable nodes)
	if isMixedFormatList(reqList) {
		self.handleMixedFormatRequirementMapping(context, reqName, reqList)
		return true
	}

	// Check if it's a list of lists (multiple assignments)
	if isListOfListsForRequirement(reqList) {
		handleMultipleRequirementAssignments(self, context, reqName, reqList)
		return true
	} else if len(reqList) == 2 {
		// For 2-element lists, we need to defer the decision until rendering
		// when all node templates are available in the namespace
		if s0, ok0 := reqList[0].(string); ok0 {
			if s1, ok1 := reqList[1].(string); ok1 {
				self.createDeferredMapping(context, reqName, s0, s1, reqList)
				return true
			}
		}
	} else if len(reqList) >= 2 && isListOfStrings(reqList) {
		// Only treat as selectable nodes if it's NOT a 2-element list
		handleSelectableNodesList(self, context, reqName, reqList)
		return true
	}
	return false
}

// =============================================================================
// SPECIFIC MAPPING HANDLERS
// =============================================================================

// handleSelectableNodesList creates mappings for selectable nodes
func handleSelectableNodesList(self *SubstitutionMappings, context *parsing.Context, reqName string, nodeList ard.List) {
	// Create multiple mappings for the count
	for i, node := range nodeList {
		nodeTemplateName := node.(string)

		// Create a unique internal key for storage
		internalKey := fmt.Sprintf("%s[%d]", reqName, i)

		// Create the mapping context
		mappingContext := context.FieldChild(internalKey, nodeTemplateName)
		mapping := NewRequirementMapping(mappingContext)

		// Use the original requirement name
		mapping.Name = reqName
		// IMPORTANT: Create separate copy of the string to avoid pointer sharing
		nodeTemplateNameCopy := nodeTemplateName
		mapping.NodeTemplateName = &nodeTemplateNameCopy
		// For selectable nodes, there's no specific requirement name - it points directly to the node
		mapping.RequirementName = nil
		mapping.isCompactForm = true // Treat as compact form for validation

		// Store using the internal key to avoid collisions
		if self.RequirementMappings == nil {
			self.RequirementMappings = make(RequirementMappings)
		}
		self.RequirementMappings[internalKey] = mapping
	}
}

// handleMultipleRequirementAssignments creates mappings for multiple requirement assignments
func handleMultipleRequirementAssignments(self *SubstitutionMappings, context *parsing.Context, reqName string, assignments ard.List) {
	// For multiple assignments, we need to create multiple mappings with the SAME requirement name
	// but different internal keys for storage in the map
	for i, assignment := range assignments {
		if assignmentList, ok := assignment.(ard.List); ok && len(assignmentList) == 2 {
			nodeTemplateName := assignmentList[0].(string)
			requirementName := assignmentList[1].(string)

			// Create a unique internal key for storage, but keep the original requirement name
			internalKey := fmt.Sprintf("%s[%d]", reqName, i)

			// Create the mapping context
			mappingContext := context.FieldChild(internalKey, assignmentList)
			mapping := NewRequirementMapping(mappingContext)

			// IMPORTANT: Use the original requirement name, not the internal key
			mapping.Name = reqName // This is what gets validated against the node type
			// IMPORTANT: Create separate copies of the strings to avoid pointer sharing
			nodeTemplateNameCopy := nodeTemplateName
			requirementNameCopy := requirementName
			mapping.NodeTemplateName = &nodeTemplateNameCopy
			mapping.RequirementName = &requirementNameCopy
			mapping.isCompactForm = false

			// Store using the internal key to avoid collisions
			if self.RequirementMappings == nil {
				self.RequirementMappings = make(RequirementMappings)
			}
			self.RequirementMappings[internalKey] = mapping
		}
	}
}

// handleTraditionalRequirementMapping creates a mapping for traditional format
func handleTraditionalRequirementMapping(self *SubstitutionMappings, context *parsing.Context, reqName string, reqList ard.List) {
	// Validate that we have exactly 2 elements and they are strings
	if len(reqList) != 2 {
		context.ReportValueMalformed("requirement mapping", "traditional format requires exactly 2 elements")
		return
	}

	nodeTemplateName, ok1 := reqList[0].(string)
	requirementName, ok2 := reqList[1].(string)
	if !ok1 || !ok2 {
		context.ReportValueMalformed("requirement mapping", "traditional format elements must be strings")
		return
	}

	// Create the mapping context
	mappingContext := context.FieldChild(reqName, reqList)
	mapping := NewRequirementMapping(mappingContext)

	// Set the mapping properties
	mapping.Name = reqName
	// IMPORTANT: Create separate copies of the strings to avoid pointer sharing
	nodeTemplateNameCopy := nodeTemplateName
	requirementNameCopy := requirementName
	mapping.NodeTemplateName = &nodeTemplateNameCopy
	mapping.RequirementName = &requirementNameCopy
	mapping.isCompactForm = false

	// Store the mapping
	if self.RequirementMappings == nil {
		self.RequirementMappings = make(RequirementMappings)
	}
	self.RequirementMappings[reqName] = mapping
}

// handleMixedFormatRequirementMapping handles lists with both traditional mappings and selectable nodes
func (self *SubstitutionMappings) handleMixedFormatRequirementMapping(context *parsing.Context, reqName string, mixedList ard.List) {
	// Process each element in the mixed list
	for i, item := range mixedList {
		if itemList, ok := item.(ard.List); ok && len(itemList) == 2 {
			// Traditional mapping format: [ node_template, requirement ]
			if nodeTemplateName, ok1 := itemList[0].(string); ok1 {
				if requirementName, ok2 := itemList[1].(string); ok2 {
					// Create a unique internal key for storage
					internalKey := fmt.Sprintf("%s[%d]", reqName, i)

					// Create the mapping context
					mappingContext := context.FieldChild(internalKey, itemList)
					mapping := NewRequirementMapping(mappingContext)

					// Set mapping properties for traditional format
					mapping.Name = reqName
					// IMPORTANT: Create separate copies of the strings to avoid pointer sharing
					nodeTemplateNameCopy := nodeTemplateName
					requirementNameCopy := requirementName
					mapping.NodeTemplateName = &nodeTemplateNameCopy
					mapping.RequirementName = &requirementNameCopy
					mapping.isCompactForm = false

					// Store the mapping
					if self.RequirementMappings == nil {
						self.RequirementMappings = make(RequirementMappings)
					}
					self.RequirementMappings[internalKey] = mapping
				}
			}
		} else if nodeTemplateName, ok := item.(string); ok {
			// Selectable node format: string
			// Create a unique internal key for storage
			internalKey := fmt.Sprintf("%s[%d]", reqName, i)

			// Create the mapping context
			mappingContext := context.FieldChild(internalKey, nodeTemplateName)
			mapping := NewRequirementMapping(mappingContext)

			// Set mapping properties for selectable node
			mapping.Name = reqName
			// IMPORTANT: Create separate copy of the string to avoid pointer sharing
			nodeTemplateNameCopy := nodeTemplateName
			mapping.NodeTemplateName = &nodeTemplateNameCopy
			mapping.RequirementName = nil // Selectable nodes don't have specific requirement names
			mapping.isCompactForm = true  // Selectable nodes are compact form

			// Store the mapping
			if self.RequirementMappings == nil {
				self.RequirementMappings = make(RequirementMappings)
			}
			self.RequirementMappings[internalKey] = mapping
		} else {
			// Invalid format in mixed list
			context.ReportValueMalformed("requirement mapping",
				fmt.Sprintf("mixed format element %d must be either [node_template, requirement] or selectable_node_name", i))
		}
	}
}

// handleMixedFormatRequirementMappingWithCounter handles lists with both traditional mappings and selectable nodes using counters
func (self *SubstitutionMappings) handleMixedFormatRequirementMappingWithCounter(context *parsing.Context, reqName string, mixedList ard.List, requirementCounters map[string]int) {
	// Process each element in the mixed list
	for _, item := range mixedList {
		if itemList, ok := item.(ard.List); ok && len(itemList) == 2 {
			// Traditional mapping format: [ node_template, requirement ]
			if nodeTemplateName, ok1 := itemList[0].(string); ok1 {
				if requirementName, ok2 := itemList[1].(string); ok2 {
					// Get next sequential index for this requirement name
					currentIndex := requirementCounters[reqName]
					requirementCounters[reqName]++

					// Create a unique internal key for storage
					internalKey := fmt.Sprintf("%s[%d]", reqName, currentIndex)

					// Create the mapping context
					mappingContext := context.FieldChild(internalKey, itemList)
					mapping := NewRequirementMapping(mappingContext)

					// Set mapping properties for traditional format
					mapping.Name = reqName
					// IMPORTANT: Create separate copies of the strings to avoid pointer sharing
					nodeTemplateNameCopy := nodeTemplateName
					requirementNameCopy := requirementName
					mapping.NodeTemplateName = &nodeTemplateNameCopy
					mapping.RequirementName = &requirementNameCopy
					mapping.isCompactForm = false

					// Store the mapping
					if self.RequirementMappings == nil {
						self.RequirementMappings = make(RequirementMappings)
					}
					self.RequirementMappings[internalKey] = mapping
				}
			}
		} else if nodeTemplateName, ok := item.(string); ok {
			// Selectable node format: string
			// Get next sequential index for this requirement name
			currentIndex := requirementCounters[reqName]
			requirementCounters[reqName]++

			// Create a unique internal key for storage
			internalKey := fmt.Sprintf("%s[%d]", reqName, currentIndex)

			// Create the mapping context
			mappingContext := context.FieldChild(internalKey, nodeTemplateName)
			mapping := NewRequirementMapping(mappingContext)

			// Set mapping properties for selectable node
			mapping.Name = reqName
			// IMPORTANT: Create separate copy of the string to avoid pointer sharing
			nodeTemplateNameCopy := nodeTemplateName
			mapping.NodeTemplateName = &nodeTemplateNameCopy
			mapping.RequirementName = nil // Selectable nodes don't have specific requirement names
			mapping.isCompactForm = true  // Selectable nodes are compact form

			// Store the mapping
			if self.RequirementMappings == nil {
				self.RequirementMappings = make(RequirementMappings)
			}
			self.RequirementMappings[internalKey] = mapping
		} else {
			// Invalid format in mixed list
			context.ReportValueMalformed("requirement mapping",
				"mixed format element must be either [node_template, requirement] or selectable_node_name")
		}
	}
}

// handleSelectableNodesListWithCounter creates mappings for selectable nodes using counters
func (self *SubstitutionMappings) handleSelectableNodesListWithCounter(context *parsing.Context, reqName string, nodeList ard.List, requirementCounters map[string]int) {
	// Create multiple mappings for the count
	for _, node := range nodeList {
		nodeTemplateName := node.(string)

		// Get next sequential index for this requirement name
		currentIndex := requirementCounters[reqName]
		requirementCounters[reqName]++

		// Create a unique internal key for storage
		internalKey := fmt.Sprintf("%s[%d]", reqName, currentIndex)

		// Create the mapping context
		mappingContext := context.FieldChild(internalKey, nodeTemplateName)
		mapping := NewRequirementMapping(mappingContext)

		// Use the original requirement name
		mapping.Name = reqName
		// IMPORTANT: Create separate copy of the string to avoid pointer sharing
		nodeTemplateNameCopy := nodeTemplateName
		mapping.NodeTemplateName = &nodeTemplateNameCopy
		// For selectable nodes, there's no specific requirement name - it points directly to the node
		mapping.RequirementName = nil
		mapping.isCompactForm = true // Treat as compact form for validation

		// Store using the internal key to avoid collisions
		if self.RequirementMappings == nil {
			self.RequirementMappings = make(RequirementMappings)
		}
		self.RequirementMappings[internalKey] = mapping
	}
}

// handleMultipleRequirementAssignmentsWithCounter creates mappings for multiple requirement assignments using counters
func (self *SubstitutionMappings) handleMultipleRequirementAssignmentsWithCounter(context *parsing.Context, reqName string, assignments ard.List, requirementCounters map[string]int) {
	// For multiple assignments, we need to create multiple mappings with the SAME requirement name
	// but different internal keys for storage in the map
	for _, assignment := range assignments {
		if assignmentList, ok := assignment.(ard.List); ok && len(assignmentList) == 2 {
			nodeTemplateName := assignmentList[0].(string)
			requirementName := assignmentList[1].(string)

			// Get next sequential index for this requirement name
			currentIndex := requirementCounters[reqName]
			requirementCounters[reqName]++

			// Create a unique internal key for storage, but keep the original requirement name
			internalKey := fmt.Sprintf("%s[%d]", reqName, currentIndex)

			// Create the mapping context
			mappingContext := context.FieldChild(internalKey, assignmentList)
			mapping := NewRequirementMapping(mappingContext)

			// IMPORTANT: Use the original requirement name, not the internal key
			mapping.Name = reqName // This is what gets validated against the node type
			// IMPORTANT: Create separate copies of the strings to avoid pointer sharing
			nodeTemplateNameCopy := nodeTemplateName
			requirementNameCopy := requirementName
			mapping.NodeTemplateName = &nodeTemplateNameCopy
			mapping.RequirementName = &requirementNameCopy
			mapping.isCompactForm = false

			// Store using the internal key to avoid collisions
			if self.RequirementMappings == nil {
				self.RequirementMappings = make(RequirementMappings)
			}
			self.RequirementMappings[internalKey] = mapping
		}
	}
}

// =============================================================================
// COMPACT FORM AND MIXED SYNTAX PROCESSING
// =============================================================================

// handleCompactFormWithCount handles compact form with count: "[ service, 2 ]: [ software2, service ]"
func (self *SubstitutionMappings) handleCompactFormWithCount(context *parsing.Context, key string, value ard.List) bool {
	// Check if the key represents compact form with count: "- database\n- 2"
	lines := strings.Split(key, "\n")
	if len(lines) != 2 {
		return false
	}

	// Parse the first line to extract requirement name
	line1 := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(line1, "- ") {
		return false
	}
	reqName := strings.TrimPrefix(line1, "- ")

	// Parse the second line to extract count
	line2 := strings.TrimSpace(lines[1])
	if !strings.HasPrefix(line2, "- ") {
		return false
	}
	countStr := strings.TrimPrefix(line2, "- ")

	// Handle both numeric count and UNBOUNDED
	var count int
	var isUnbounded bool
	if countStr == "UNBOUNDED" {
		isUnbounded = true
		count = -1 // Use -1 as marker for UNBOUNDED
	} else if parsedCount, err := strconv.Atoi(countStr); err == nil && parsedCount > 0 {
		count = parsedCount
	} else {
		return false
	}

	// Validate the value is a 2-element list
	if len(value) != 2 {
		context.ReportValueMalformed("requirement mapping",
			"compact form with count requires exactly 2 values: [node_template, requirement]")
		return false
	}

	nodeTemplateName, ok1 := value[0].(string)
	requirementName, ok2 := value[1].(string)
	if !ok1 || !ok2 {
		context.ReportValueMalformed("requirement mapping",
			"compact form with count values must be strings")
		return false
	}

	// For UNBOUNDED, validate against the substituted node type (TOSCA v2.0 compliant)
	if isUnbounded {
		// TOSCA v2.0: UNBOUNDED mappings should validate against the substituted node type,
		// not the target node template. Defer validation until rendering when NodeType is available.
		count = 1 // Use temporary count, will be corrected during rendering
	}

	// Create multiple mappings for the count
	for i := 0; i < count; i++ {
		internalKey := fmt.Sprintf("%s[%d]", reqName, i)

		// Create the mapping context with the original value, not individual elements
		mappingContext := context.FieldChild(internalKey, value)
		mapping := NewRequirementMapping(mappingContext)

		// Use the original requirement name
		mapping.Name = reqName
		mapping.NodeTemplateName = &nodeTemplateName
		mapping.RequirementName = &requirementName
		mapping.isCompactForm = false // This is explicit form with node and requirement

		// Mark UNBOUNDED mappings for later processing
		if isUnbounded {
			mapping.isUnbounded = true
		}

		// Store using the internal key to avoid collisions
		if self.RequirementMappings == nil {
			self.RequirementMappings = make(RequirementMappings)
		}
		self.RequirementMappings[internalKey] = mapping
	}

	return true
}

// processMixedRequirementMapping handles all types of requirement mapping syntaxes in a unified way
func (self *SubstitutionMappings) processMixedRequirementMapping(context *parsing.Context, reqKey interface{}, reqValue interface{}, itemIndex int) {
	switch key := reqKey.(type) {
	case string:
		// Standard format: "service: [ software1, service ]"
		if self.isSerializedListKey(key) {
			// Parse the serialized list: "[ service, 2 ]" or "[ service, UNBOUNDED ]"
			self.processCompactFormRequirementMapping(context, key, reqValue, itemIndex)
		} else {
			// Regular string key - treat as standard format
			self.processStandardRequirementMapping(context, key, reqValue, itemIndex)
		}
	default:
		// Handle YAMLKey type (complex keys like [ service, 2 ])
		if keyStr := self.extractKeyString(reqKey); keyStr != "" {
			if self.isSerializedListKey(keyStr) {
				// Parse the serialized list: "[ service, 2 ]" or "[ service, UNBOUNDED ]"
				self.processCompactFormRequirementMapping(context, keyStr, reqValue, itemIndex)
			} else {
				// Regular string key - treat as standard format
				self.processStandardRequirementMapping(context, keyStr, reqValue, itemIndex)
			}
		} else {
			context.ReportValueMalformed("requirement mapping",
				fmt.Sprintf("unsupported key type for requirement mapping: %T", reqKey))
		}
	}
}

// processMixedRequirementMappingWithCounter handles all types of requirement mapping syntaxes with proper counter management
func (self *SubstitutionMappings) processMixedRequirementMappingWithCounter(context *parsing.Context, reqKey interface{}, reqValue interface{}, itemIndex int, requirementCounters map[string]int) {
	switch key := reqKey.(type) {
	case string:
		// Standard format: "service: [ software1, service ]"
		if self.isSerializedListKey(key) {
			// Parse the serialized list: "[ service, 2 ]" or "[ service, UNBOUNDED ]"
			self.processCompactFormRequirementMappingWithCounter(context, key, reqValue, itemIndex, requirementCounters)
		} else {
			// Regular string key - treat as standard format
			self.processStandardRequirementMappingWithCounter(context, key, reqValue, itemIndex, requirementCounters)
		}
	default:
		// Handle YAMLKey type (complex keys like [ service, 2 ])
		if keyStr := self.extractKeyString(reqKey); keyStr != "" {
			if self.isSerializedListKey(keyStr) {
				// Parse the serialized list: "[ service, 2 ]" or "[ service, UNBOUNDED ]"
				self.processCompactFormRequirementMappingWithCounter(context, keyStr, reqValue, itemIndex, requirementCounters)
			} else {
				// Regular string key - treat as standard format
				self.processStandardRequirementMappingWithCounter(context, keyStr, reqValue, itemIndex, requirementCounters)
			}
		} else {
			context.ReportValueMalformed("requirement mapping",
				fmt.Sprintf("unsupported key type for requirement mapping: %T", reqKey))
		}
	}
}

// processStandardRequirementMapping handles standard format: "service: [ software1, service ]"
func (self *SubstitutionMappings) processStandardRequirementMapping(context *parsing.Context, reqName string, reqValue interface{}, itemIndex int) {
	// Handle single string value (selectable node)
	if reqStr, ok := reqValue.(string); ok {
		// Single node template name - this is a selectable node mapping
		internalKey := fmt.Sprintf("%s[%d]", reqName, itemIndex)
		mappingContext := context.FieldChild(internalKey, reqValue)
		mapping := NewRequirementMapping(mappingContext)

		mapping.Name = reqName
		mapping.NodeTemplateName = &reqStr
		mapping.RequirementName = nil // For selectable nodes, points directly to the node
		mapping.isCompactForm = true  // Selectable nodes are compact form

		if self.RequirementMappings == nil {
			self.RequirementMappings = make(RequirementMappings)
		}
		self.RequirementMappings[internalKey] = mapping
		return
	}

	if reqList, ok := reqValue.(ard.List); ok {
		// Check for mixed format first
		if isMixedFormatList(reqList) {
			self.handleMixedFormatRequirementMapping(context, reqName, reqList)
			return
		}

		if len(reqList) == 2 {
			// Traditional format: [ node_template, requirement ]
			nodeTemplateName, ok1 := reqList[0].(string)
			requirementName, ok2 := reqList[1].(string)
			if ok1 && ok2 {
				// Create the mapping
				internalKey := fmt.Sprintf("%s[%d]", reqName, itemIndex)
				mappingContext := context.FieldChild(internalKey, reqValue)
				mapping := NewRequirementMapping(mappingContext)

				mapping.Name = reqName
				mapping.NodeTemplateName = &nodeTemplateName
				mapping.RequirementName = &requirementName
				mapping.isCompactForm = false

				if self.RequirementMappings == nil {
					self.RequirementMappings = make(RequirementMappings)
				}
				self.RequirementMappings[internalKey] = mapping
				return
			}
		}

		// Could be selectable nodes list
		if isListOfStrings(reqList) {
			handleSelectableNodesList(self, context, reqName, reqList)
			return
		}

		// Could be multiple assignments list
		if isListOfStringPairs(reqList) {
			handleMultipleRequirementAssignments(self, context, reqName, reqList)
			return
		}
	}

	context.ReportValueMalformed("requirement mapping",
		fmt.Sprintf("unsupported value format for requirement '%s'", reqName))
}

// processStandardRequirementMappingWithCounter handles standard format with counter: "service: [ software1, service ]"
func (self *SubstitutionMappings) processStandardRequirementMappingWithCounter(context *parsing.Context, reqName string, reqValue interface{}, itemIndex int, requirementCounters map[string]int) {
	// Handle single string value (selectable node)
	if reqStr, ok := reqValue.(string); ok {
		// Get next sequential index for this requirement name
		currentIndex := requirementCounters[reqName]
		requirementCounters[reqName]++

		// Single node template name - this is a selectable node mapping
		internalKey := fmt.Sprintf("%s[%d]", reqName, currentIndex)
		mappingContext := context.FieldChild(internalKey, reqValue)
		mapping := NewRequirementMapping(mappingContext)

		mapping.Name = reqName
		// IMPORTANT: Create separate copy of the string to avoid pointer sharing
		reqStrCopy := reqStr
		mapping.NodeTemplateName = &reqStrCopy
		mapping.RequirementName = nil // For selectable nodes, points directly to the node
		mapping.isCompactForm = true  // Selectable nodes are compact form

		if self.RequirementMappings == nil {
			self.RequirementMappings = make(RequirementMappings)
		}
		self.RequirementMappings[internalKey] = mapping
		return
	}

	if reqList, ok := reqValue.(ard.List); ok {
		// Check for mixed format first
		if isMixedFormatList(reqList) {
			self.handleMixedFormatRequirementMappingWithCounter(context, reqName, reqList, requirementCounters)
			return
		}

		if len(reqList) == 2 {
			// Traditional format: [ node_template, requirement ]
			nodeTemplateName, ok1 := reqList[0].(string)
			requirementName, ok2 := reqList[1].(string)
			if ok1 && ok2 {
				// Get next sequential index for this requirement name
				currentIndex := requirementCounters[reqName]
				requirementCounters[reqName]++

				// Create the mapping
				internalKey := fmt.Sprintf("%s[%d]", reqName, currentIndex)
				mappingContext := context.FieldChild(internalKey, reqValue)
				mapping := NewRequirementMapping(mappingContext)

				mapping.Name = reqName
				// IMPORTANT: Create separate copies of the strings to avoid pointer sharing
				nodeTemplateNameCopy := nodeTemplateName
				requirementNameCopy := requirementName
				mapping.NodeTemplateName = &nodeTemplateNameCopy
				mapping.RequirementName = &requirementNameCopy
				mapping.isCompactForm = false

				if self.RequirementMappings == nil {
					self.RequirementMappings = make(RequirementMappings)
				}
				self.RequirementMappings[internalKey] = mapping
				return
			}
		}

		// Could be selectable nodes list
		if isListOfStrings(reqList) {
			self.handleSelectableNodesListWithCounter(context, reqName, reqList, requirementCounters)
			return
		}

		// Could be multiple assignments list
		if isListOfStringPairs(reqList) {
			self.handleMultipleRequirementAssignmentsWithCounter(context, reqName, reqList, requirementCounters)
			return
		}
	}

	context.ReportValueMalformed("requirement mapping",
		fmt.Sprintf("unsupported value format for requirement '%s'", reqName))
}

// processCompactFormRequirementMapping handles compact form: "[ service, 2 ]: [ software2, service ]"
func (self *SubstitutionMappings) processCompactFormRequirementMapping(context *parsing.Context, keyStr string, reqValue interface{}, itemIndex int) {
	// Parse the serialized list key
	reqName, countStr, ok := self.parseSerializedListKey(keyStr)
	if !ok {
		context.ReportValueMalformed("requirement mapping",
			fmt.Sprintf("invalid compact form key format: %s", keyStr))
		return
	}

	// Validate the value is a 2-element list
	reqList, ok := reqValue.(ard.List)
	if !ok || len(reqList) != 2 {
		context.ReportValueMalformed("requirement mapping",
			"compact form requires exactly 2 values: [node_template, requirement]")
		return
	}

	nodeTemplateName, ok1 := reqList[0].(string)
	requirementName, ok2 := reqList[1].(string)
	if !ok1 || !ok2 {
		context.ReportValueMalformed("requirement mapping",
			"compact form values must be strings")
		return
	}

	// Handle count
	var count int
	var isUnbounded bool
	if countStr == "UNBOUNDED" {
		isUnbounded = true
		count = 1 // Use temporary count, will be corrected during rendering
	} else if parsedCount, err := strconv.Atoi(countStr); err == nil && parsedCount > 0 {
		count = parsedCount
	} else {
		context.ReportValueMalformed("requirement mapping",
			fmt.Sprintf("invalid count value: %s", countStr))
		return
	}

	// Create mappings for the count
	for i := 0; i < count; i++ {
		internalKey := fmt.Sprintf("%s[%d]", reqName, itemIndex*1000+i) // Use itemIndex*1000 to ensure uniqueness

		mappingContext := context.FieldChild(internalKey, reqValue)
		mapping := NewRequirementMapping(mappingContext)

		mapping.Name = reqName
		mapping.NodeTemplateName = &nodeTemplateName
		mapping.RequirementName = &requirementName
		mapping.isCompactForm = false

		if isUnbounded {
			mapping.isUnbounded = true
		}

		// Set the requirement index
		mapping.RequirementIndex = &i

		if self.RequirementMappings == nil {
			self.RequirementMappings = make(RequirementMappings)
		}
		self.RequirementMappings[internalKey] = mapping
	}
}

// processCompactFormRequirementMappingWithCounter handles compact form with counter: "[ service, 2 ]: [ software2, service ]"
func (self *SubstitutionMappings) processCompactFormRequirementMappingWithCounter(context *parsing.Context, keyStr string, reqValue interface{}, itemIndex int, requirementCounters map[string]int) {
	// Parse the serialized list key
	reqName, countStr, ok := self.parseSerializedListKey(keyStr)
	if !ok {
		context.ReportValueMalformed("requirement mapping",
			fmt.Sprintf("invalid compact form key format: %s", keyStr))
		return
	}

	// Validate the value is a 2-element list
	reqList, ok := reqValue.(ard.List)
	if !ok || len(reqList) != 2 {
		context.ReportValueMalformed("requirement mapping",
			"compact form requires exactly 2 values: [node_template, requirement]")
		return
	}

	nodeTemplateName, ok1 := reqList[0].(string)
	requirementName, ok2 := reqList[1].(string)
	if !ok1 || !ok2 {
		context.ReportValueMalformed("requirement mapping",
			"compact form values must be strings")
		return
	}

	// Handle count
	var count int
	var isUnbounded bool
	if countStr == "UNBOUNDED" {
		isUnbounded = true
		count = 1 // Use temporary count, will be corrected during rendering
	} else if parsedCount, err := strconv.Atoi(countStr); err == nil && parsedCount > 0 {
		count = parsedCount
	} else {
		context.ReportValueMalformed("requirement mapping",
			fmt.Sprintf("invalid count value: %s", countStr))
		return
	}

	// Get the starting index for this requirement name
	startIndex := requirementCounters[reqName]

	// Create mappings for the count
	for i := 0; i < count; i++ {
		// Calculate the current index
		currentIndex := startIndex + i

		internalKey := fmt.Sprintf("%s[%d]", reqName, currentIndex)

		mappingContext := context.FieldChild(internalKey, reqValue)
		mapping := NewRequirementMapping(mappingContext)

		mapping.Name = reqName
		// IMPORTANT: Create separate copies of the strings to avoid pointer sharing
		nodeTemplateNameCopy := nodeTemplateName
		requirementNameCopy := requirementName
		mapping.NodeTemplateName = &nodeTemplateNameCopy
		mapping.RequirementName = &requirementNameCopy
		mapping.isCompactForm = false

		if isUnbounded {
			mapping.isUnbounded = true
		}

		// Set the requirement index using the currentIndex, not the loop index
		currentIndexCopy := currentIndex
		mapping.RequirementIndex = &currentIndexCopy

		if self.RequirementMappings == nil {
			self.RequirementMappings = make(RequirementMappings)
		}
		self.RequirementMappings[internalKey] = mapping
	}

	// Update the counter to reflect all mappings created
	requirementCounters[reqName] = startIndex + count
}

// =============================================================================
// CLEANUP AND UTILITY FUNCTIONS
// =============================================================================

// cleanupProcessedRequirements removes processed requirements from the context data
func (self *SubstitutionMappings) cleanupProcessedRequirements(context *parsing.Context, reqMap ard.Map, processedRequirements map[string]bool) {
	// Create a new map without the processed requirements
	newReqMap := make(ard.Map)
	for reqName, reqValue := range reqMap {
		if !processedRequirements[reqName.(string)] {
			newReqMap[reqName] = reqValue
		}
	}

	// Update the context data with the cleaned map
	if mapData, ok := context.Data.(ard.Map); ok {
		newMap := make(ard.Map)
		for key, value := range mapData {
			if key == "requirements" {
				newMap[key] = newReqMap
			} else {
				newMap[key] = value
			}
		}
		context.Data = newMap
	}
}

// clearOriginalRequirements clears the original requirements list from context
func (self *SubstitutionMappings) clearOriginalRequirements(context *parsing.Context) {
	if mapData, ok := context.Data.(ard.Map); ok {
		newMap := make(ard.Map)
		for key, value := range mapData {
			if key != "requirements" {
				newMap[key] = value
			}
		}
		context.Data = newMap
	}
}

// createDeferredMapping creates a mapping that will be resolved during rendering
func (self *SubstitutionMappings) createDeferredMapping(context *parsing.Context, reqName, nodeTemplateName, requirementName string, reqList ard.List) {
	// Create a special mapping that will be resolved during rendering
	mappingContext := context.FieldChild(reqName, reqList)
	mapping := NewRequirementMapping(mappingContext)

	// Set the mapping properties with a special marker
	mapping.Name = reqName
	mapping.NodeTemplateName = &nodeTemplateName
	mapping.RequirementName = &requirementName
	mapping.isCompactForm = false // Default to false, will be updated during rendering

	// Add a special marker to indicate this needs resolution
	mapping.RequirementIndex = nil // Use nil to indicate needs resolution

	// Store the mapping
	if self.RequirementMappings == nil {
		self.RequirementMappings = make(RequirementMappings)
	}
	self.RequirementMappings[reqName] = mapping
}

// =============================================================================
// POST-PROCESSING FUNCTIONS
// =============================================================================

func (self *SubstitutionMappings) postProcessCompactFormWithCount(context *parsing.Context) {
	// Look for requirement mappings with the special COMPACT_FORM_WITH_COUNT and COMPACT_FORM_WITH_INDEX markers
	toRemove := make([]string, 0)
	toAdd := make(map[string]*RequirementMapping)

	for key, mapping := range self.RequirementMappings {
		if strings.HasPrefix(mapping.Name, "COMPACT_FORM_WITH_COUNT:") {
			// Parse the marker: "COMPACT_FORM_WITH_COUNT:reqName:count:nodeTemplate:requirement"
			parts := strings.Split(mapping.Name, ":")
			if len(parts) == 5 {
				reqName := parts[1]
				countStr := parts[2]
				nodeTemplateName := parts[3]
				requirementName := parts[4]

				var count int
				if countStr == "UNBOUNDED" {
					// TOSCA v2.0: UNBOUNDED mappings should validate against the substituted node type
					// Defer validation until rendering when NodeType is available.
					count = 1 // Use temporary count, will be corrected during rendering
				} else if parsedCount, err := strconv.Atoi(countStr); err == nil && parsedCount > 0 {
					count = parsedCount
				} else {
					continue // Skip invalid count
				}

				// Create multiple mappings for the count
				for i := 0; i < count; i++ {
					internalKey := fmt.Sprintf("%s[%d]", reqName, i)

					// Create the mapping context
					mappingContext := context.FieldChild(internalKey, mapping.Context.Data)
					newMapping := NewRequirementMapping(mappingContext)

					// Use the original requirement name
					newMapping.Name = reqName
					newMapping.NodeTemplateName = &nodeTemplateName
					newMapping.RequirementName = &requirementName
					newMapping.isCompactForm = false // This is explicit form with node and requirement
					// Set the requirement index for validation
					newMapping.RequirementIndex = &i

					// Mark UNBOUNDED mappings for later processing
					if countStr == "UNBOUNDED" {
						newMapping.isUnbounded = true
					}

					// Store using the internal key to avoid collisions
					toAdd[internalKey] = newMapping
				}

				// Mark the original mapping for removal
				toRemove = append(toRemove, key)
			}
		}
	}

	// Remove the marker mappings and add the new ones
	for _, key := range toRemove {
		delete(self.RequirementMappings, key)
	}

	for key, mapping := range toAdd {
		self.RequirementMappings[key] = mapping
	}
}

// =============================================================================
// QUERY FUNCTIONS
// =============================================================================

func (self *SubstitutionMappings) IsRequirementMapped(nodeTemplate *NodeTemplate, requirementName string) bool {
	for _, mapping := range self.RequirementMappings {
		if mapping.NodeTemplate == nodeTemplate {
			if (mapping.RequirementName != nil) && (*mapping.RequirementName == requirementName) {
				return true
			}
		}
	}
	return false
}

// =============================================================================
// MAIN RENDERING FUNCTIONS
// =============================================================================

func (self *SubstitutionMappings) Render(inputDefinitions ParameterDefinitions) {
	if self.NodeType == nil {
		return
	}

	self.renderCapabilityMappings()
	self.renderRequirementMappings()
	self.renderPropertyMappings(inputDefinitions)
	self.renderAttributeMappings()
	self.renderInterfaceMappings()

	// Add validation for get_attribute functions in outputs
	self.validateOutputsGetAttributeFunctions()
}

// =============================================================================
// OUTPUT VALIDATION FUNCTIONS
// =============================================================================

func (self *SubstitutionMappings) validateOutputsGetAttributeFunctions() {
	serviceContext := self.findServiceTemplateContext()
	if serviceContext == nil {
		return
	}

	outputsContext, ok := serviceContext.GetFieldChild("outputs")
	if !ok {
		return
	}

	// Traverse all outputs and validate get_attribute functions
	for _, outputContext := range outputsContext.FieldChildren() {
		if valueContext, ok := outputContext.GetFieldChild("value"); ok {
			self.validateGetAttributeInValue(valueContext)
		}
	}
}

func (self *SubstitutionMappings) validateGetAttributeInValue(context *parsing.Context) {
	if context.Is(ard.TypeList) {
		if list, ok := context.Data.(ard.List); ok {
			for i, item := range list {
				itemContext := context.ListChild(i, item)
				self.validateGetAttributeInValue(itemContext)
			}
		}
	} else if context.Is(ard.TypeMap) {
		if map_, ok := context.Data.(ard.Map); ok {
			// Check for get_attribute function
			if getAttribute, ok := map_["$get_attribute"]; ok {
				self.validateGetAttributeFunction(context, getAttribute)
			}

			// Recursively check other map entries
			for key, value := range map_ {
				if key != "$get_attribute" {
					childContext := context.MapChild(key, value)
					self.validateGetAttributeInValue(childContext)
				}
			}
		}
	}
}

func (self *SubstitutionMappings) validateGetAttributeFunction(context *parsing.Context, functionData interface{}) {
	if args, ok := functionData.(ard.List); ok && len(args) >= 4 {
		nodeTemplateNameArg, ok1 := args[0].(string)
		requirementNameArg, ok2 := args[1].(string)
		attributeNameArg, ok3 := args[2].(string)

		var relationshipIndex int = 0
		if len(args) >= 4 {
			if idx, ok4 := args[3].(int); ok4 {
				relationshipIndex = idx
			} else if idxFloat, ok4 := args[3].(float64); ok4 {
				relationshipIndex = int(idxFloat)
			}
		}

		if ok1 && ok2 && ok3 {
			// Validate the relationship attribute and index
			self.validateGetAttributeRelationshipAccess(context, nodeTemplateNameArg, requirementNameArg, attributeNameArg, relationshipIndex)
		}
	}
}

func (self *SubstitutionMappings) validateGetAttributeRelationshipAccess(context *parsing.Context, nodeTemplateName, requirementName, attributeName string, relationshipIndex int) {
	// Find the service template context
	serviceContext := self.findServiceTemplateContext()
	if serviceContext == nil {
		return
	}

	nodeTemplatesContext, ok := serviceContext.GetFieldChild("node_templates")
	if !ok {
		return
	}

	nodeTemplateContext, ok := nodeTemplatesContext.GetFieldChild(nodeTemplateName)
	if !ok {
		context.ReportValueMalformed("get_attribute function",
			fmt.Sprintf("node template '%s' not found", nodeTemplateName))
		return
	}

	// Count actual relationships for this requirement
	actualRelationshipCount := self.countRequirementsInNodeTemplate(nodeTemplateContext, requirementName)

	// Validate relationship index
	if relationshipIndex >= actualRelationshipCount {
		context.ReportValueMalformed("get_attribute function",
			fmt.Sprintf("relationship index %d exceeds available relationships (%d) for requirement '%s' in node template '%s'",
				relationshipIndex, actualRelationshipCount, requirementName, nodeTemplateName))
		return
	}

	// Get node template type to find requirement definitions
	nodeTypeNameContext, ok := nodeTemplateContext.GetFieldChild("type")
	if !ok {
		return
	}

	nodeTypeName, ok := nodeTypeNameContext.Data.(string)
	if !ok {
		return
	}

	// Look up the node type
	if nodeTypePtr, ok := context.Namespace.Lookup(nodeTypeName); ok {
		if nodeType, ok := nodeTypePtr.(*NodeType); ok {
			// Check if requirement exists in node type
			reqDef, ok := nodeType.RequirementDefinitions[requirementName]
			if !ok {
				context.ReportValueMalformed("get_attribute function",
					fmt.Sprintf("requirement '%s' not found in node type '%s'", requirementName, nodeTypeName))
				return
			}

			// Get relationship type and validate attribute
			relationshipType := self.getRelationshipType(reqDef)
			if relationshipType == nil {
				context.ReportValueMalformed("get_attribute function",
					fmt.Sprintf("relationship type not found for requirement '%s'", requirementName))
				return
			}

			// Check if attribute exists in relationship type
			if _, ok := relationshipType.AttributeDefinitions[attributeName]; !ok {
				context.ReportValueMalformed("get_attribute function",
					fmt.Sprintf("attribute '%s' not found in relationship type '%s'", attributeName, relationshipType.Name))
				return
			}
		}
	}
}

// =============================================================================
// CAPABILITY MAPPING RENDERING
// =============================================================================

func (self *SubstitutionMappings) renderCapabilityMappings() {
	for name, mapping := range self.CapabilityMappings {
		if definition, ok := self.NodeType.CapabilityDefinitions[name]; ok {
			if mappedDefinition, ok := mapping.GetCapabilityDefinition(); ok {
				self.validateTypeCompatibility(definition.CapabilityType, mappedDefinition.CapabilityType)
			}
		} else {
			mapping.Context.Clone(name).ReportReferenceNotFound("capability", self.NodeType)
		}
	}
}

// =============================================================================
// REQUIREMENT MAPPING RENDERING
// =============================================================================

func (self *SubstitutionMappings) renderRequirementMappings() {
	// STEP 1: Apply TOSCA rule for automatic expansion of single mappings BEFORE rendering
	if self.NodeType != nil {
		// Group mappings by requirement name for expansion check
		preRenderGroups := make(map[string][]*RequirementMapping)
		for _, mapping := range self.RequirementMappings {
			reqName := mapping.Name
			preRenderGroups[reqName] = append(preRenderGroups[reqName], mapping)
		}

		// Apply auto-expansion for single mappings that need it
		expansionsToAdd := make(map[string]*RequirementMapping)
		for reqName, mappings := range preRenderGroups {
			if definition, ok := self.NodeType.RequirementDefinitions[reqName]; ok {
				if definition.CountRange != nil && definition.CountRange.Range != nil {
					minCount := int(definition.CountRange.Range.Lower)
					actualCount := len(mappings)

					// TOSCA rule: if there's exactly one mapping but minCount > 1, expand automatically
					if actualCount == 1 && minCount > 1 {
						originalMapping := mappings[0]

						// Set the requirement index for the original mapping (index 0)
						if originalMapping.RequirementIndex == nil {
							indexZero := 0
							originalMapping.RequirementIndex = &indexZero
						}

						// Create additional mappings to reach minCount
						for i := 1; i < minCount; i++ {
							internalKey := fmt.Sprintf("%s[%d]", reqName, i)

							// Clone the original mapping
							mappingContext := originalMapping.Context.FieldChild(internalKey, originalMapping.Context.Data)
							newMapping := NewRequirementMapping(mappingContext)

							// Copy properties from original mapping
							newMapping.Name = originalMapping.Name
							if originalMapping.NodeTemplateName != nil {
								nodeTemplateNameCopy := *originalMapping.NodeTemplateName
								newMapping.NodeTemplateName = &nodeTemplateNameCopy
							}
							if originalMapping.RequirementName != nil {
								requirementNameCopy := *originalMapping.RequirementName
								newMapping.RequirementName = &requirementNameCopy
							}
							newMapping.isCompactForm = originalMapping.isCompactForm

							// Set the requirement index
							indexCopy := i
							newMapping.RequirementIndex = &indexCopy

							expansionsToAdd[internalKey] = newMapping
						}
					}
				}
			}
		}

		// Add the expanded mappings to the main RequirementMappings
		for key, mapping := range expansionsToAdd {
			self.RequirementMappings[key] = mapping
		}
	}

	// STEP 2: Now render all mappings (including expanded ones)
	for _, mapping := range self.RequirementMappings {
		mapping.Render()
	}

	// Handle UNBOUNDED mappings: expand them to the correct count based on NodeType
	if self.NodeType != nil {
		toAdd := make(map[string]*RequirementMapping)
		toRemove := make([]string, 0)

		for key, mapping := range self.RequirementMappings {
			if mapping.isUnbounded && mapping.RequirementIndex == nil {
				// Only process original UNBOUNDED placeholders, not auto-expanded ones
				// Find the requirement definition in the substituted node type
				if reqDef, ok := self.NodeType.RequirementDefinitions[mapping.Name]; ok {
					if reqDef.CountRange != nil && reqDef.CountRange.Range != nil {
						minCount := int(reqDef.CountRange.Range.Lower)

						// Count existing mappings for this requirement
						existingCount := 0
						hasAutoExpanded := false
						for _, existingMapping := range self.RequirementMappings {
							if existingMapping.Name == mapping.Name {
								// Count all mappings except UNBOUNDED placeholders (RequirementIndex == nil)
								if !existingMapping.isUnbounded || existingMapping.RequirementIndex != nil {
									existingCount++
								}
								// Check if this mapping was created by auto-expansion (has RequirementIndex set)
								if existingMapping.isUnbounded && existingMapping.RequirementIndex != nil {
									hasAutoExpanded = true
								}
							}
						}

						// If auto-expansion already created enough mappings, just remove the original UNBOUNDED placeholder
						if hasAutoExpanded && existingCount >= minCount {
							toRemove = append(toRemove, key)
						} else if minCount > existingCount {
							// Create additional mappings to reach the minimum count
							// Remove the original UNBOUNDED mapping
							toRemove = append(toRemove, key)

							// Create additional mappings starting from existingCount
							for i := existingCount; i < minCount; i++ {
								internalKey := fmt.Sprintf("%s[%d]", mapping.Name, i)

								// Create new mapping
								mappingContext := mapping.Context.FieldChild(internalKey, mapping.Context.Data)
								newMapping := NewRequirementMapping(mappingContext)

								// Copy properties from original mapping
								newMapping.Name = mapping.Name
								// IMPORTANT: Create separate copies of the strings to avoid pointer sharing
								if mapping.NodeTemplateName != nil {
									nodeTemplateNameCopy := *mapping.NodeTemplateName
									newMapping.NodeTemplateName = &nodeTemplateNameCopy
								}
								if mapping.RequirementName != nil {
									requirementNameCopy := *mapping.RequirementName
									newMapping.RequirementName = &requirementNameCopy
								}
								newMapping.isCompactForm = mapping.isCompactForm
								newMapping.isUnbounded = true // Keep the flag for normalization

								// Set the requirement index
								indexCopy := i
								newMapping.RequirementIndex = &indexCopy

								toAdd[internalKey] = newMapping
							}
						} else {
							// If we already have enough mappings, just remove the UNBOUNDED placeholder
							toRemove = append(toRemove, key)
						}
					}
				}
			}
		}

		// Apply the changes
		for _, key := range toRemove {
			delete(self.RequirementMappings, key)
		}
		for key, mapping := range toAdd {
			self.RequirementMappings[key] = mapping
		}
	}

	// Post-process: handle deferred selectable node mappings that were converted during rendering
	toAdd := make(map[string]*RequirementMapping)
	toRemove := make([]string, 0)

	for key, mapping := range self.RequirementMappings {
		// Check if this mapping was converted to selectable nodes during rendering
		if mapping.isCompactForm && mapping.RequirementName == nil {
			// This is now a selectable node mapping
			// We need to create separate mappings for each selectable node

			// Get the original list from the context data
			if contextData, ok := mapping.Context.Data.(ard.List); ok && len(contextData) == 2 {
				if node1, ok1 := contextData[0].(string); ok1 {
					if node2, ok2 := contextData[1].(string); ok2 {
						// Create mapping for first node
						key1 := fmt.Sprintf("%s[0]", mapping.Name)
						mappingContext1 := mapping.Context.FieldChild(key1, node1)
						mapping1 := NewRequirementMapping(mappingContext1)
						mapping1.Name = mapping.Name
						// IMPORTANT: Create separate copy of the string to avoid pointer sharing
						node1Copy := node1
						mapping1.NodeTemplateName = &node1Copy
						mapping1.RequirementName = nil
						mapping1.isCompactForm = true
						toAdd[key1] = mapping1

						// Create mapping for second node
						key2 := fmt.Sprintf("%s[1]", mapping.Name)
						mappingContext2 := mapping.Context.FieldChild(key2, node2)
						mapping2 := NewRequirementMapping(mappingContext2)
						mapping2.Name = mapping.Name
						// IMPORTANT: Create separate copy of the string to avoid pointer sharing
						node2Copy := node2
						mapping2.NodeTemplateName = &node2Copy
						mapping2.RequirementName = nil
						mapping2.isCompactForm = true
						toAdd[key2] = mapping2

						// Mark original mapping for removal
						toRemove = append(toRemove, key)
					}
				}
			}
		}
	}

	// Apply the changes
	for _, key := range toRemove {
		delete(self.RequirementMappings, key)
	}
	for key, mapping := range toAdd {
		self.RequirementMappings[key] = mapping
	}

	// Group mappings by requirement name to handle multiple assignments
	requirementGroups := make(map[string][]*RequirementMapping)

	for _, mapping := range self.RequirementMappings {
		reqName := mapping.Name
		requirementGroups[reqName] = append(requirementGroups[reqName], mapping)
	}

	// Validate each requirement group
	for reqName, mappings := range requirementGroups {
		if definition, ok := self.NodeType.RequirementDefinitions[reqName]; ok {
			// Validate count if specified
			if definition.CountRange != nil && definition.CountRange.Range != nil {
				minCount := int(definition.CountRange.Range.Lower)
				actualCount := len(mappings)
				maxCount := -1
				if definition.CountRange.Range.Upper != math.MaxUint64 {
					maxCount = int(definition.CountRange.Range.Upper)
				}
				validateCountInRange(self.Context, requirementMappingStr, reqName, actualCount, minCount, maxCount)
			}
		} else {
			// Report error only once per requirement name, not for each mapping
			validateReferenceNotFound(self.Context, "requirement", self.NodeType)
		}
	}
}

// =============================================================================
// DIRECTIVE VALIDATION FUNCTIONS
// =============================================================================

func (self *SubstitutionMappings) validateSelectDirectiveForNode(nodeTemplateName string) bool {
	var nodeTemplateType *NodeTemplate
	if nodeTemplate, ok := self.Context.Namespace.LookupForType(nodeTemplateName, reflect.TypeOf(nodeTemplateType)); ok {
		nodeTemplatePtr := nodeTemplate.(*NodeTemplate)

		if nodeTemplatePtr.Directives == nil {
			return false
		}

		for _, directive := range *nodeTemplatePtr.Directives {
			if directive == "select" {
				return true
			}
		}
	}
	return false
}

// =============================================================================
// PROPERTY MAPPING RENDERING
// =============================================================================

func (self *SubstitutionMappings) renderPropertyMappings(inputDefinitions ParameterDefinitions) {
	self.PropertyMappings.Render(inputDefinitions)
	for name, mapping := range self.PropertyMappings {
		switch {
		case self.isCapabilityMapping(mapping):
			self.validateCapabilityPropertyMapping(mapping)
		case self.isRelationshipMapping(mapping):
			self.validateRelationshipPropertyMapping(mapping)
		default:
			self.validateDirectPropertyMapping(name, mapping)
		}
	}
}

// =============================================================================
// PROPERTY MAPPING VALIDATION FUNCTIONS
// =============================================================================

func (self *SubstitutionMappings) isCapabilityMapping(mapping *PropertyMapping) bool {
	return mapping.CapabilityName != nil && mapping.CapabilityPropertyName != nil
}

func (self *SubstitutionMappings) isRelationshipMapping(mapping *PropertyMapping) bool {
	return mapping.RequirementName != nil && mapping.RelationshipPropertyName != nil
}

func (self *SubstitutionMappings) validateCapabilityPropertyMapping(mapping *PropertyMapping) {
	capability, ok := self.NodeType.CapabilityDefinitions[*mapping.CapabilityName]
	if !ok {
		mapping.Context.ReportReferenceNotFound("capability", self.NodeType)
		return
	}
	if capability.CapabilityType == nil {
		return
	}
	propertyDef, ok := capability.CapabilityType.PropertyDefinitions[*mapping.CapabilityPropertyName]
	if !ok {
		mapping.Context.ReportReferenceNotFound("capability property", capability.CapabilityType)
		return
	}
	propertyDef.Render()
	if mapping.InputDefinition != nil {
		self.validateTypeCompatibility(propertyDef.DataType, mapping.InputDefinition.DataType)
	}
}

func (self *SubstitutionMappings) validateRelationshipPropertyMapping(mapping *PropertyMapping) {
	reqDef, ok := self.NodeType.RequirementDefinitions[*mapping.RequirementName]
	if !ok {
		mapping.Context.ReportReferenceNotFound("requirement", self.NodeType)
		return
	}
	if self.shouldValidateRelationshipIndex(mapping) {
		if !self.validateRelationshipIndex(mapping) {
			return
		}
	}
	relationshipType := self.getRelationshipType(reqDef)
	if relationshipType == nil {
		mapping.Context.ReportReferenceNotFound("relationship type", reqDef)
		return
	}
	propertyDef, ok := relationshipType.PropertyDefinitions[*mapping.RelationshipPropertyName]
	if !ok {
		mapping.Context.ReportReferenceNotFound("relationship property", relationshipType)
		return
	}
	propertyDef.Render()
	if mapping.InputDefinition != nil {
		self.validateRelationshipPropertyMappingTypes(mapping, propertyDef)
	}
}

func (self *SubstitutionMappings) validateRelationshipPropertyMappingTypes(mapping *PropertyMapping, propertyDef *PropertyDefinition) {
	if self.isAllRelationshipMapping(mapping) {
		self.validateAllRelationshipMapping(mapping, propertyDef)
	} else {
		self.validateSingleRelationshipMapping(mapping, propertyDef)
	}
}

func (self *SubstitutionMappings) validateDirectPropertyMapping(name string, mapping *PropertyMapping) {
	definition, ok := self.NodeType.PropertyDefinitions[name]
	if !ok {
		mapping.Context.Clone(name).ReportReferenceNotFound("property", self.NodeType)
		return
	}
	definition.Render()
	if mapping.InputDefinition != nil {
		self.validateInputMapping(definition, mapping)
	} else if mapping.Value != nil {
		self.validatePropertyMapping(definition, mapping)
	}
}

func (self *SubstitutionMappings) validateInputMapping(definition *PropertyDefinition, mapping *PropertyMapping) {
	if definition.DataType != nil && mapping.InputDefinition.DataType != nil {
		self.validateTypeCompatibility(definition.DataType, mapping.InputDefinition.DataType)
	}
}

func (self *SubstitutionMappings) validatePropertyMapping(definition *PropertyDefinition, mapping *PropertyMapping) {
	if definition.DataType == nil {
		return
	}

	if mapping.Value.DataType != nil {
		self.validateTypeCompatibility(definition.DataType, mapping.Value.DataType)
	} else {
		mapping.Value.RenderProperty(definition.DataType, definition)
	}
}

func (self *SubstitutionMappings) validateTypeCompatibility(expected, actual parsing.EntityPtr) {
	if expected != nil && actual != nil {
		if !self.Context.Hierarchy.IsCompatible(expected, actual) {
			self.Context.ReportIncompatibleType(expected, actual)
		}
	}
}

// =============================================================================
// ATTRIBUTE MAPPING RENDERING
// =============================================================================

func (self *SubstitutionMappings) renderAttributeMappings() {
	self.AttributeMappings.EnsureRender()
	for name, mapping := range self.AttributeMappings {
		if mapping.IsCapabilityMapping() {
			self.validateCapabilityAttributeMapping(mapping)
		} else if mapping.IsRelationshipMapping() {
			self.validateRelationshipAttributeMapping(mapping)
		} else {
			// Direct attribute mapping validation
			if definition, ok := self.NodeType.AttributeDefinitions[name]; ok {
				if definition.DataType != nil && mapping.Attribute != nil && mapping.Attribute.DataType != nil {
					self.validateTypeCompatibility(definition.DataType, mapping.Attribute.DataType)
				}
			} else {
				mapping.Context.Clone(name).ReportReferenceNotFound("attribute", self.NodeType)
			}
		}

		// Always validate output mapping if present (TOSCA 2.0)
		if mapping.IsOutputMapping() {
			self.validateOutputMapping(mapping)
		}
	}
}

// =============================================================================
// ATTRIBUTE MAPPING VALIDATION FUNCTIONS
// =============================================================================

func (self *SubstitutionMappings) validateCapabilityAttributeMapping(mapping *AttributeMapping) {
	capability, ok := self.NodeType.CapabilityDefinitions[*mapping.CapabilityName]
	if !ok {
		mapping.Context.ReportReferenceNotFound("capability", self.NodeType)
		return
	}

	if capability.CapabilityType == nil {
		return
	}

	attributeDef, ok := capability.CapabilityType.AttributeDefinitions[*mapping.CapabilityAttributeName]
	if !ok {
		mapping.Context.ReportReferenceNotFound("capability attribute", capability.CapabilityType)
		return
	}

	attributeDef.Render()
}

func (self *SubstitutionMappings) validateRelationshipAttributeMapping(mapping *AttributeMapping) {
	reqDef, ok := self.NodeType.RequirementDefinitions[*mapping.RequirementName]
	if !ok {
		mapping.Context.ReportReferenceNotFound("requirement", self.NodeType)
		return
	}

	// Add relationship index validation for attribute mappings
	if self.shouldValidateAttributeRelationshipIndex(mapping) {
		if !self.validateAttributeRelationshipIndex(mapping) {
			return
		}
	}

	relationshipType := self.getRelationshipType(reqDef)
	if relationshipType == nil {
		mapping.Context.ReportReferenceNotFound("relationship type", reqDef)
		return
	}

	attributeDef, ok := relationshipType.AttributeDefinitions[*mapping.RelationshipAttributeName]
	if !ok {
		mapping.Context.ReportReferenceNotFound("relationship attribute", relationshipType)
		return
	}

	attributeDef.Render()
}

func (self *SubstitutionMappings) shouldValidateAttributeRelationshipIndex(mapping *AttributeMapping) bool {
	return mapping.RelationshipIndex != nil && *mapping.RelationshipIndex != -1
}

func (self *SubstitutionMappings) validateAttributeRelationshipIndex(mapping *AttributeMapping) bool {
	actualCount := self.countActualRequirements(*mapping.RequirementName)

	if *mapping.RelationshipIndex >= actualCount {
		mapping.Context.ReportValueMalformed("relationship index",
			fmt.Sprintf("index %d exceeds actual requirement instances %d for requirement '%s'",
				*mapping.RelationshipIndex, actualCount, *mapping.RequirementName))
		return false
	}
	return true
}

func (self *SubstitutionMappings) validateOutputMapping(mapping *AttributeMapping) {
	// Find the service template context to validate outputs
	serviceContext := self.findServiceTemplateContext()
	if serviceContext == nil {
		return
	}

	outputsContext, ok := serviceContext.GetFieldChild("outputs")
	if !ok {
		mapping.Context.ReportValueMalformed("attribute mapping", "service template has no outputs section")
		return
	}

	if _, ok := outputsContext.GetFieldChild(*mapping.OutputName); !ok {
		mapping.Context.ReportValueMalformed("attribute mapping",
			fmt.Sprintf("output '%s' not found in service template", *mapping.OutputName))
	}
}

// =============================================================================
// INTERFACE MAPPING RENDERING
// =============================================================================

func (self *SubstitutionMappings) renderInterfaceMappings() {
	for name, mapping := range self.InterfaceMappings {
		self.validateInterfaceMapping(name, mapping)
	}
}

func (self *SubstitutionMappings) validateInterfaceMapping(name string, mapping *InterfaceMapping) {
	definition, ok := self.NodeType.InterfaceDefinitions[name]
	if !ok {
		mapping.Context.Clone(name).ReportReferenceNotFound("interface", self.NodeType)
		return
	}

	// Check if this is a TOSCA 1.3 interface mapping
	if _, isTosca13 := mapping.OperationMappings["__TOSCA_1_3__"]; isTosca13 {
		// Skip validation for TOSCA 1.3 format - it uses different semantics
		return
	}

	// TOSCA 2.0: Validate that each operation in the mapping exists in the interface definition
	if definition.InterfaceType != nil {
		for operationName, workflowName := range mapping.OperationMappings {
			if _, ok := definition.InterfaceType.OperationDefinitions[operationName]; !ok {
				mapping.Context.MapChild(operationName, workflowName).ReportReferenceNotFound("operation", definition.InterfaceType)
			}
		}
	}
}

// =============================================================================
// NORMALIZATION FUNCTIONS
// =============================================================================

// normalizeCapabilities normalizes capability mappings
func (self *SubstitutionMappings) normalizeCapabilities(normalSubstitution *normal.Substitution, normalServiceTemplate *normal.ServiceTemplate) {
	for _, mapping := range self.CapabilityMappings {
		if mapping.NodeTemplate != nil && mapping.NodeTemplateCapabilityName != nil {
			if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]; ok {
				normalSubstitution.CapabilityPointers[mapping.Name] = normalNodeTemplate.NewPointer(*mapping.NodeTemplateCapabilityName)
			}
		}
	}
}

// normalizeRequirements normalizes requirement mappings
func (self *SubstitutionMappings) normalizeRequirements(normalSubstitution *normal.Substitution, normalServiceTemplate *normal.ServiceTemplate) {
	// Group mappings by original requirement name
	requirementGroups := make(map[string][]*RequirementMapping)
	for _, mapping := range self.RequirementMappings {
		reqName := mapping.Name
		requirementGroups[reqName] = append(requirementGroups[reqName], mapping)
	}

	for reqName, mappings := range requirementGroups {
		if len(mappings) == 1 {
			self.normalizeSingleRequirementMapping(mappings[0], normalSubstitution, normalServiceTemplate)
		} else {
			self.normalizeMultipleRequirementMappings(reqName, mappings, normalSubstitution, normalServiceTemplate)
		}
	}
}

// normalizeProperties normalizes property mappings
func (self *SubstitutionMappings) normalizeProperties(normalSubstitution *normal.Substitution, normalServiceTemplate *normal.ServiceTemplate) {
	for _, mapping := range self.PropertyMappings {
		cleanKey := self.generateMappingKey(mapping)
		if mapping.NodeTemplate != nil && mapping.PropertyName != nil {
			if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]; ok {
				normalSubstitution.PropertyPointers[cleanKey] = normalNodeTemplate.NewPointer(*mapping.PropertyName)
			}
		} else if mapping.Value != nil {
			normalSubstitution.PropertyValues[cleanKey] = mapping.Value.Normalize()
		} else if mapping.InputName != nil {
			normalSubstitution.InputPointers[cleanKey] = normal.NewPointer(*mapping.InputName)
		}
	}
}

// normalizeAttributes normalizes attribute mappings
func (self *SubstitutionMappings) normalizeAttributes(normalSubstitution *normal.Substitution, normalServiceTemplate *normal.ServiceTemplate) {
	for _, mapping := range self.AttributeMappings {
		var key string
		if mapping.IsCapabilityMapping() || mapping.IsRelationshipMapping() {
			key = self.generateAttributeMappingKey(mapping)
		} else {
			key = mapping.Name
		}
		if mapping.OutputName != nil {
			normalSubstitution.OutputPointers[key] = normal.NewPointer(*mapping.OutputName)
		}
		if mapping.IsCapabilityMapping() {
			self.normalizeCapabilityAttributeMapping(mapping, normalSubstitution, normalServiceTemplate, key)
		} else if mapping.IsRelationshipMapping() {
			self.normalizeRelationshipAttributeMapping(mapping, normalSubstitution, normalServiceTemplate, key)
		}
		if mapping.NodeTemplate != nil && mapping.AttributeName != nil {
			if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]; ok {
				normalSubstitution.AttributePointers[key] = normalNodeTemplate.NewPointer(*mapping.AttributeName)
			}
		}
	}
}

// normalizeInterfaces normalizes interface mappings
func (self *SubstitutionMappings) normalizeInterfaces(normalSubstitution *normal.Substitution, normalServiceTemplate *normal.ServiceTemplate) {
	// TOSCA 2.0: Interface mappings map operations to workflows
	// TOSCA 1.3: Interface mappings map interfaces to node templates (stored with special key)
	for interfaceName, mapping := range self.InterfaceMappings {
		// Check for TOSCA 1.3 format first
		if tosca13Info, ok := mapping.OperationMappings["__TOSCA_1_3__"]; ok {
			// Parse the special format: "nodeTemplate:interfaceName"
			parts := strings.Split(tosca13Info, ":")
			if len(parts) == 2 {
				nodeTemplateName := parts[0]
				targetInterfaceName := parts[1]

				// Check if the node template exists
				if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[nodeTemplateName]; ok {
					// Create interface pointer that points to the node template's interface
					normalSubstitution.InterfacePointers[interfaceName] = normalNodeTemplate.NewPointer(targetInterfaceName)
				}
			}
		} else {
			// Standard TOSCA 2.0 processing
			for operationName, workflowName := range mapping.OperationMappings {
				// Create a unique key for the interface operation mapping
				key := fmt.Sprintf("%s.%s", interfaceName, operationName)

				// Check if the workflow exists in the service template
				if normalServiceTemplate.Workflows != nil {
					if _, ok := normalServiceTemplate.Workflows[workflowName]; ok {
						// Create a pointer to the workflow for interface mapping
						// This creates a connection between the interface operation and the workflow
						normalSubstitution.InterfacePointers[key] = normal.NewPointer(workflowName)
					}
				}
			}
		}
	}
}

// Normalize refactored for modularity
func (self *SubstitutionMappings) Normalize(normalServiceTemplate *normal.ServiceTemplate) *normal.Substitution {
	if self.NodeType == nil {
		return nil
	}
	normalSubstitution := normalServiceTemplate.NewSubstitution()
	normalSubstitution.Type = parsing.GetCanonicalName(self.NodeType)
	if metadata, ok := self.NodeType.GetMetadata(); ok {
		normalSubstitution.TypeMetadata = metadata
	}
	self.normalizeCapabilities(normalSubstitution, normalServiceTemplate)
	self.normalizeRequirements(normalSubstitution, normalServiceTemplate)
	self.normalizeProperties(normalSubstitution, normalServiceTemplate)
	self.normalizeAttributes(normalSubstitution, normalServiceTemplate)
	self.normalizeInterfaces(normalSubstitution, normalServiceTemplate)
	return normalSubstitution
}

// =============================================================================
// UTILITY AND HELPER FUNCTIONS
// =============================================================================

// findServiceTemplateContext returns the parent context named "service_template"
func (self *SubstitutionMappings) findServiceTemplateContext() *parsing.Context {
	context := self.Context
	for context.Parent != nil && context.Name != "service_template" {
		context = context.Parent
	}
	if context.Name == "service_template" {
		return context
	}
	return nil
}

// countRequirementsInTargetNode counts requirementName in nodeTemplateName
func (self *SubstitutionMappings) countRequirementsInTargetNode(context *parsing.Context, nodeTemplateName, requirementName string) int {
	serviceContext := self.findServiceTemplateContext()
	if serviceContext == nil {
		return 0
	}
	nodeTemplatesContext, ok := serviceContext.GetFieldChild("node_templates")
	if !ok {
		return 0
	}
	nodeTemplateContext, ok := nodeTemplatesContext.GetFieldChild(nodeTemplateName)
	if !ok {
		return 0
	}
	return self.countRequirementsInNodeTemplate(nodeTemplateContext, requirementName)
}

// countRequirementsInNodeTemplate counts requirementName in a nodeTemplate context
func (self *SubstitutionMappings) countRequirementsInNodeTemplate(nodeTemplateContext *parsing.Context, requirementName string) int {
	count := 0
	requirementsContext, ok := nodeTemplateContext.GetFieldChild("requirements")
	if !ok || !requirementsContext.ValidateType(ard.TypeList) {
		return count
	}
	for _, reqData := range requirementsContext.Data.(ard.List) {
		if reqMap, ok := reqData.(ard.Map); ok {
			if _, exists := reqMap[requirementName]; exists {
				count++
			}
		}
	}
	return count
}

// getRelationshipType returns the RelationshipType from a RequirementDefinition
func (self *SubstitutionMappings) getRelationshipType(reqDef *RequirementDefinition) *RelationshipType {
	if reqDef.RelationshipDefinition != nil && reqDef.RelationshipDefinition.RelationshipType != nil {
		return reqDef.RelationshipDefinition.RelationshipType
	}
	return nil
}

// shouldValidateRelationshipIndex returns true if the mapping's index should be validated
func (self *SubstitutionMappings) shouldValidateRelationshipIndex(mapping *PropertyMapping) bool {
	return mapping.RelationshipIndex != nil && *mapping.RelationshipIndex != -1
}

// validateRelationshipIndex checks if the mapping's index is within allowed range
func (self *SubstitutionMappings) validateRelationshipIndex(mapping *PropertyMapping) bool {
	if reqDef, ok := self.NodeType.RequirementDefinitions[*mapping.RequirementName]; ok {
		if reqDef.CountRange == nil || reqDef.CountRange.Range == nil {
			return true
		}
		if reqDef.CountRange.Range.Upper == math.MaxUint64 {
			return true
		}
		maxCount := int(reqDef.CountRange.Range.Upper)
		if *mapping.RelationshipIndex >= maxCount {
			mapping.Context.ReportValueMalformed("relationship index",
				fmt.Sprintf("index %d exceeds maximum requirement count %d for requirement '%s' in substituted type '%s'",
					*mapping.RelationshipIndex, maxCount, *mapping.RequirementName, self.NodeType.Name))
			return false
		}
	}
	return true
}

// isAllRelationshipMapping returns true if mapping is for all relationships
func (self *SubstitutionMappings) isAllRelationshipMapping(mapping *PropertyMapping) bool {
	return mapping.RelationshipIndex != nil && *mapping.RelationshipIndex == -1
}

// validateAllRelationshipMapping checks input type for ALL relationship mapping
func (self *SubstitutionMappings) validateAllRelationshipMapping(mapping *PropertyMapping, propertyDef *PropertyDefinition) {
	if mapping.InputDefinition.DataType == nil {
		return
	}
	if mapping.InputDefinition.DataType.Name != "list" {
		self.Context.ReportValueMalformed("property mapping",
			fmt.Sprintf("input for ALL relationship mapping must be a list type, got %s",
				mapping.InputDefinition.DataType.Name))
		return
	}
	self.validateListEntryType(mapping, propertyDef)
}

// validateSingleRelationshipMapping checks type compatibility for a single mapping
func (self *SubstitutionMappings) validateSingleRelationshipMapping(mapping *PropertyMapping, propertyDef *PropertyDefinition) {
	if propertyDef.DataType != nil && mapping.InputDefinition.DataType != nil {
		self.validateTypeCompatibility(propertyDef.DataType, mapping.InputDefinition.DataType)
	}
}

// countActualRequirements counts all occurrences of requirementName in all node templates
func (self *SubstitutionMappings) countActualRequirements(requirementName string) int {
	count := 0
	serviceContext := self.findServiceTemplateContext()
	if serviceContext == nil {
		return count
	}
	nodeTemplatesContext, ok := serviceContext.GetFieldChild("node_templates")
	if !ok {
		return count
	}
	for _, nodeTemplateContext := range nodeTemplatesContext.FieldChildren() {
		count += self.countRequirementsInNodeTemplate(nodeTemplateContext, requirementName)
	}
	return count
}

// =============================================================================
// VALIDATION HELPER FUNCTIONS
// =============================================================================

// validateCountInRange checks that actual is between min and max (if max > 0)
func validateCountInRange(context *parsing.Context, what, name string, actual, min, max int) {
	if actual < min {
		context.ReportValueMalformed(what,
			fmt.Sprintf("'%s' has %d mappings but requires at least %d", name, actual, min))
	}
	if max > 0 && actual > max {
		context.ReportValueMalformed(what,
			fmt.Sprintf("'%s' has %d mappings but allows maximum %d", name, actual, max))
	}
}

// validateReferenceNotFound reports a missing reference
func validateReferenceNotFound(context *parsing.Context, what string, entity interface{}) {
	context.ReportReferenceNotFound(what, entity)
}

// validateListEntryType checks type compatibility for list entries
func (self *SubstitutionMappings) validateListEntryType(mapping *PropertyMapping, propertyDef *PropertyDefinition) {
	entrySchema := mapping.InputDefinition.DataType.EntrySchema
	if entrySchema == nil || entrySchema.DataType == nil {
		return
	}
	if propertyDef.DataType != nil {
		self.validateTypeCompatibility(propertyDef.DataType, entrySchema.DataType)
	}
}

// normalizeSingleRequirementMapping normalizes a single requirement mapping
func (self *SubstitutionMappings) normalizeSingleRequirementMapping(mapping *RequirementMapping, normalSubstitution *normal.Substitution, normalServiceTemplate *normal.ServiceTemplate) {
	if mapping.NodeTemplate == nil {
		return
	}
	normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[mapping.NodeTemplate.Name]
	if !ok {
		return
	}
	if mapping.RequirementName == nil {
		normalSubstitution.RequirementPointer[mapping.Name] = normalNodeTemplate.NewPointer("")
	} else {
		normalSubstitution.RequirementPointer[mapping.Name] = normalNodeTemplate.NewPointer(*mapping.RequirementName)
	}
}

// normalizeMultipleRequirementMappings normalizes multiple requirement mappings
func (self *SubstitutionMappings) normalizeMultipleRequirementMappings(reqName string, mappings []*RequirementMapping, normalSubstitution *normal.Substitution, normalServiceTemplate *normal.ServiceTemplate) {
	for i, mapping := range mappings {
		if mapping.NodeTemplate == nil {
			continue
		}

		nodeTemplateName := mapping.NodeTemplate.Name

		normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[nodeTemplateName]
		if !ok {
			continue
		}

		indexedName := fmt.Sprintf("%s[%d]", reqName, i)

		if mapping.RequirementName == nil {
			normalSubstitution.RequirementPointer[indexedName] = normalNodeTemplate.NewPointer("")
		} else {
			requirementName := *mapping.RequirementName
			normalSubstitution.RequirementPointer[indexedName] = normalNodeTemplate.NewPointer(requirementName)
		}
	}
}

// generateMappingKey generates a key for property mappings
func (self *SubstitutionMappings) generateMappingKey(mapping *PropertyMapping) string {
	if mapping.CapabilityName != nil && mapping.CapabilityPropertyName != nil {
		return fmt.Sprintf("[CAPABILITY, %s, %s]", *mapping.CapabilityName, *mapping.CapabilityPropertyName)
	}
	if mapping.RequirementName != nil && mapping.RelationshipPropertyName != nil {
		indexStr := self.getIndexString(mapping.RelationshipIndex)
		return fmt.Sprintf("[RELATIONSHIP, %s, %s, %s]", *mapping.RequirementName, indexStr, *mapping.RelationshipPropertyName)
	}
	return mapping.Name
}

// generateAttributeMappingKey generates a key for attribute mappings
func (self *SubstitutionMappings) generateAttributeMappingKey(mapping *AttributeMapping) string {
	if mapping.CapabilityName != nil && mapping.CapabilityAttributeName != nil {
		return fmt.Sprintf("[CAPABILITY, %s, %s]", *mapping.CapabilityName, *mapping.CapabilityAttributeName)
	}
	if mapping.RequirementName != nil && mapping.RelationshipAttributeName != nil {
		indexStr := self.getIndexString(mapping.RelationshipIndex)
		return fmt.Sprintf("[RELATIONSHIP, %s, %s, %s]", *mapping.RequirementName, indexStr, *mapping.RelationshipAttributeName)
	}
	return mapping.Name
}

// normalizeCapabilityAttributeMapping validates and normalizes capability attribute mapping
func (self *SubstitutionMappings) normalizeCapabilityAttributeMapping(mapping *AttributeMapping, normalSubstitution *normal.Substitution, normalServiceTemplate *normal.ServiceTemplate, key string) {
	capability, ok := self.NodeType.CapabilityDefinitions[*mapping.CapabilityName]
	if !ok || capability.CapabilityType == nil {
		return
	}
	if _, ok := capability.CapabilityType.AttributeDefinitions[*mapping.CapabilityAttributeName]; !ok {
		return
	}
	if mapping.OutputName != nil {
		normalSubstitution.OutputPointers[key] = normal.NewPointer(*mapping.OutputName)
	}
}

// normalizeRelationshipAttributeMapping validates and normalizes relationship attribute mapping
func (self *SubstitutionMappings) normalizeRelationshipAttributeMapping(mapping *AttributeMapping, normalSubstitution *normal.Substitution, normalServiceTemplate *normal.ServiceTemplate, key string) {
	reqDef, ok := self.NodeType.RequirementDefinitions[*mapping.RequirementName]
	if !ok {
		return
	}
	relationshipType := self.getRelationshipType(reqDef)
	if relationshipType == nil {
		return
	}
	if _, ok := relationshipType.AttributeDefinitions[*mapping.RelationshipAttributeName]; !ok {
		return
	}
	if mapping.OutputName != nil {
		normalSubstitution.OutputPointers[key] = normal.NewPointer(*mapping.OutputName)
	}
}

// getIndexString returns the string representation of a relationship index
func (self *SubstitutionMappings) getIndexString(relationshipIndex *int) string {
	if relationshipIndex == nil {
		return "0"
	}
	if *relationshipIndex == -1 {
		return "ALL"
	}
	return fmt.Sprintf("%d", *relationshipIndex)
}
