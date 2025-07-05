package tosca_v2_0

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// RequirementMapping
//
// Attaches to SubstitutionMappings
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.8
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.8
// [TOSCA-Simple-Profile-YAML-v1.1] @ 2.10, 2.11
// [TOSCA-Simple-Profile-YAML-v1.0] @ 2.10, 2.11
//

type RequirementMapping struct {
	*Entity `name:"requirement mapping"`
	Name    string

	NodeTemplateName *string
	RequirementName  *string

	NodeTemplate *NodeTemplate          `traverse:"ignore" json:"-" yaml:"-"`
	Requirement  *RequirementAssignment `traverse:"ignore" json:"-" yaml:"-"`

	// Flag to track if this is compact form for validation
	isCompactForm bool
	// Track which index of the requirement when there are multiple instances
	RequirementIndex *int
	// Flag to track if this is UNBOUNDED mapping requiring deferred validation
	isUnbounded bool
}

func NewRequirementMapping(context *parsing.Context) *RequirementMapping {
	return &RequirementMapping{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadRequirementMapping(context *parsing.Context) parsing.EntityPtr {
	self := NewRequirementMapping(context)

	// Check if this is a compact form with count (multiline string key)
	if strings.Contains(context.Name, "\n") {
		// Check if the key represents compact form with count: "- database\n- 2"
		lines := strings.Split(context.Name, "\n")

		if len(lines) == 2 {
			// Parse the first line to extract requirement name
			line1 := strings.TrimSpace(lines[0])
			if strings.HasPrefix(line1, "- ") {
				reqName := strings.TrimPrefix(line1, "- ")

				// Parse the second line to extract count
				line2 := strings.TrimSpace(lines[1])
				if strings.HasPrefix(line2, "- ") {
					countStr := strings.TrimPrefix(line2, "- ")

					// Handle both numeric count and UNBOUNDED
					var count int
					var isUnbounded bool
					if countStr == "UNBOUNDED" {
						isUnbounded = true
						count = -1 // Use -1 as marker for UNBOUNDED
					} else if parsedCount, err := strconv.Atoi(countStr); err == nil && parsedCount > 0 {
						// This is a count - generates multiple mappings
						count = parsedCount
					} else {
						// Invalid count format
						return self
					}

					// Validate the value is a 2-element list
					if context.ValidateType(ard.TypeList) {
						if list, ok := context.Data.(ard.List); ok && len(list) == 2 {
							if nodeTemplateName, ok1 := list[0].(string); ok1 {
								if requirementName, ok2 := list[1].(string); ok2 {
									// Return a special marker to indicate this should be handled by substitution mappings
									if isUnbounded {
										self.Name = fmt.Sprintf("COMPACT_FORM_WITH_COUNT:%s:UNBOUNDED:%s:%s", reqName, nodeTemplateName, requirementName)
									} else {
										self.Name = fmt.Sprintf("COMPACT_FORM_WITH_COUNT:%s:%d:%s:%s", reqName, count, nodeTemplateName, requirementName)
									}
									return self
								}
							}
						}
					}
				}
			}
		}
	}

	// Support four forms:
	// 1. Compact form (string): requirement_name: node_template_name
	// 2. List form (2 strings): requirement_name: [node_template_name, requirement_name]
	// 3. Multiple assignments (list of lists): requirement_name: [[node1, req1], [node2, req2]]
	// 4. List of selectable nodes: requirement_name: [selectable_node1, selectable_node2]

	if context.ValidateType(ard.TypeString) {
		// Format 1: Compact form (selectable node)
		nodeTemplateName := context.Data.(string)
		self.NodeTemplateName = &nodeTemplateName
		self.RequirementName = nil // For selectable nodes, RequirementName is nil
		self.isCompactForm = true
	} else if context.ValidateType(ard.TypeList) {
		list := context.Data.(ard.List)

		if len(list) == 2 {
			// Check if it's a list of selectable nodes or normal list form
			if self.isListOfSelectableNodes(context, list) {
				// This should be handled by substitution mappings
				context.ReportValueMalformed("requirement mapping",
					"list of selectable nodes should be processed by substitution mappings parser")
			} else if strings := context.ReadStringListFixed(2); strings != nil {
				// Format 2: Simple list form
				self.NodeTemplateName = &(*strings)[0]
				self.RequirementName = &(*strings)[1]
				self.isCompactForm = false
			} else {
				context.ReportValueMalformed("requirement mapping", "2-element list must contain strings")
			}
		} else if len(list) > 0 {
			// Check if it's a list of lists (multiple assignments) or list of selectable nodes
			if self.isListOfLists(list) {
				// This should have been handled by substitution mappings, but if not, report error
				context.ReportValueMalformed("requirement mapping",
					"list of lists format should be processed by substitution mappings parser")
			} else if self.isListOfSelectableNodes(context, list) {
				// This should be handled by substitution mappings
				context.ReportValueMalformed("requirement mapping",
					"list of selectable nodes should be processed by substitution mappings parser")
			} else {
				context.ReportValueMalformed("requirement mapping", "list elements must be strings or lists")
			}
		} else {
			context.ReportValueMalformed("requirement mapping", "empty list not allowed")
		}
	} else {
		context.ReportValueMalformed("requirement mapping", "must be either a string (compact form) or a list")
	}

	return self
}

// Check if this is a list of lists (multiple assignments)
func (self *RequirementMapping) isListOfLists(list ard.List) bool {
	for _, item := range list {
		if _, ok := item.(ard.List); !ok {
			return false
		}
	}
	return true
}

// Check if this is a list of selectable nodes
func (self *RequirementMapping) isListOfSelectableNodes(context *parsing.Context, list ard.List) bool {
	// All elements must be strings (node template names)
	for _, item := range list {
		if _, ok := item.(string); !ok {
			return false
		}
	}

	// Check if all referenced nodes have select directive
	allHaveSelect := true
	for _, item := range list {
		nodeTemplateName := item.(string)
		if !self.nodeHasSelectDirective(context, nodeTemplateName) {
			allHaveSelect = false
		}
	}

	return allHaveSelect
}

// Check if a specific node template has select directive
func (self *RequirementMapping) nodeHasSelectDirective(context *parsing.Context, nodeTemplateName string) bool {
	var nodeTemplateType *NodeTemplate
	if nodeTemplate, ok := context.Namespace.LookupForType(nodeTemplateName, reflect.TypeOf(nodeTemplateType)); ok {
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

// ([parsing.Mappable] interface)
func (self *RequirementMapping) GetKey() string {
	return self.Name
}

func (self *RequirementMapping) GetRequirementDefinition() (*RequirementDefinition, bool) {
	if (self.Requirement != nil) && (self.NodeTemplate != nil) {
		return self.Requirement.GetDefinition(self.NodeTemplate)
	} else {
		return nil, false
	}
}

// ([parsing.Renderable] interface)
func (self *RequirementMapping) Render() {
	// Avoid rendering more than once
	self.renderOnce.Do(self.render)
}

func (self *RequirementMapping) render() {
	logRender.Debug("requirement mapping")

	if self.NodeTemplateName == nil {
		return
	}

	nodeTemplateName := *self.NodeTemplateName
	var nodeTemplateType *NodeTemplate
	if nodeTemplate, ok := self.Context.Namespace.LookupForType(nodeTemplateName, reflect.TypeOf(nodeTemplateType)); ok {
		self.NodeTemplate = nodeTemplate.(*NodeTemplate)

		self.NodeTemplate.Render()

		// DEFERRED RESOLUTION: Check if this is a 2-element mapping that needs to be resolved
		// This happens when we have a mapping like [db1, db2] that could be either:
		// 1. [node_template, requirement] pair, or
		// 2. List of selectable nodes
		if self.RequirementName != nil && self.RequirementIndex == nil {
			requirementName := *self.RequirementName

			// Check if the second element (requirementName) is actually a node template with select directive
			var secondNodeType *NodeTemplate
			if secondNode, ok := self.Context.Namespace.LookupForType(requirementName, reflect.TypeOf(secondNodeType)); ok {
				secondNodePtr := secondNode.(*NodeTemplate)

				// Check if both nodes have select directive
				hasSelect1 := false
				hasSelect2 := false

				if self.NodeTemplate.Directives != nil {
					for _, directive := range *self.NodeTemplate.Directives {
						if directive == "select" {
							hasSelect1 = true
							break
						}
					}
				}

				if secondNodePtr.Directives != nil {
					for _, directive := range *secondNodePtr.Directives {
						if directive == "select" {
							hasSelect2 = true
							break
						}
					}
				}

				if hasSelect1 && hasSelect2 {
					// Both are selectable nodes - convert to selectable node mapping
					self.RequirementName = nil // Points directly to the nodes
					self.isCompactForm = true  // Selectable nodes are compact form
				} else {
					// Not both selectable - keep as [node_template, requirement] pair
					self.isCompactForm = false
				}
			} else {
				// Second element is not a node template - must be a requirement name
				self.isCompactForm = false
			}
		}

		// CRITICAL: For compact form, REQUIRE that the target node has select directive
		// This validation must happen BEFORE checking quirks because it's a TOSCA 2.0 requirement

		if self.isCompactForm {
			hasSelect := self.hasSelectDirective()

			if !hasSelect {
				// This is a TOSCA 2.0 requirement, so it must be enforced regardless of quirks
				self.Context.ReportValueMalformed("requirement mapping",
					fmt.Sprintf("compact form 'database: %s' requires target node to have 'directives: [select]' but node has no select directive", nodeTemplateName))
				return // Stop processing if validation fails
			} else {
				// For selectable nodes, set RequirementName to nil as it points directly to the node
				self.RequirementName = nil
			}
		}

		// ONLY validate requirement existence if RequirementName is not nil
		// For selectable nodes, RequirementName is nil because it points directly to the node
		if self.RequirementName != nil {
			name := *self.RequirementName

			// If RequirementIndex is set, use it to directly access the requirement
			if self.RequirementIndex != nil {
				index := *self.RequirementIndex

				// Count matching requirements to validate index
				matchingRequirements := make([]*RequirementAssignment, 0)
				for _, requirement := range self.NodeTemplate.Requirements {
					if requirement.Name == name {
						matchingRequirements = append(matchingRequirements, requirement)
					}
				}

				if len(matchingRequirements) > index {
					self.Requirement = matchingRequirements[index]
				} else {
					// TOSCA v2.0: Check for dangling requirements in node_type
					// A requirement defined in the node_type but not explicitly assigned in the node_template
					// is considered "dangling" and available for substitution mappings
					danglingRequirement := self.findDanglingRequirement(name)
					if danglingRequirement != nil {
						// Create a placeholder requirement assignment for the dangling requirement
						// This represents the implicit assignment that will be created by the processor
						self.Requirement = nil // Will be handled during normalization
					} else {
						// TOSCA v2.0: For UNBOUNDED mappings, this is allowed because implicit assignments
						// will be created by the TOSCA processor when needed
						if self.isUnbounded {
							// For UNBOUNDED mappings, create a placeholder requirement assignment
							// This represents an implicit assignment that will be created by the processor
							self.Requirement = nil // Will be handled during normalization
						} else {
							self.Context.ListChild(1, name).ReportReferenceNotFound("requirement", self.NodeTemplate)
						}
					}
				}
			} else {
				// Original logic for single requirement
				found := false
				for _, requirement := range self.NodeTemplate.Requirements {
					if requirement.Name == name {
						if found {
							self.Context.ListChild(1, name).ReportReferenceAmbiguous("requirement", self.NodeTemplate)
							break
						} else {
							self.Requirement = requirement
							found = true
						}
					}
				}

				if !found {
					// TOSCA v2.0: Check for dangling requirements in node_type
					// A requirement defined in the node_type but not explicitly assigned in the node_template
					// is considered "dangling" and available for substitution mappings
					danglingRequirement := self.findDanglingRequirement(name)
					if danglingRequirement != nil {
						// Create a placeholder requirement assignment for the dangling requirement
						// This represents the implicit assignment that will be created by the processor
						self.Requirement = nil // Will be handled during normalization
					} else {
						// TOSCA v2.0: For UNBOUNDED mappings, this is allowed because implicit assignments
						// will be created by the TOSCA processor when needed
						if self.isUnbounded {
							// For UNBOUNDED mappings, create a placeholder requirement assignment
							// This represents an implicit assignment that will be created by the processor
							self.Requirement = nil // Will be handled during normalization
						} else {
							self.Context.ListChild(1, name).ReportReferenceNotFound("requirement", self.NodeTemplate)
						}
					}
				}
			}
		}

		// Now check quirks for other validations (but requirement existence was already validated above)
		if self.Context.HasQuirk(parsing.QuirkSubstitutionMappingsRequirementsPermissive) {
			return
		}

	} else {
		self.Context.ListChild(0, nodeTemplateName).ReportUnknown("node template")
		return
	}
}

// Check if the target node template has select directive
func (self *RequirementMapping) hasSelectDirective() bool {
	if self.NodeTemplate == nil {
		return false
	}

	if self.NodeTemplate.Directives == nil {
		return false
	}

	for _, directive := range *self.NodeTemplate.Directives {
		if directive == "select" {
			return true
		}
	}

	return false
}

// findDanglingRequirement checks if a requirement is defined in the node_type but not explicitly assigned in the node_template
// This is a TOSCA v2.0 concept for substitution mappings - dangling requirements are available for mapping
func (self *RequirementMapping) findDanglingRequirement(requirementName string) *RequirementDefinition {
	// Check if we have a valid node template and its node type
	if self.NodeTemplate == nil || self.NodeTemplate.NodeType == nil {
		return nil
	}

	// Look for the requirement definition in the node type
	if reqDef, ok := self.NodeTemplate.NodeType.RequirementDefinitions[requirementName]; ok {
		// Check if this requirement is NOT explicitly assigned in the node template
		// (i.e., it's dangling and available for substitution mappings)
		isExplicitlyAssigned := false
		for _, req := range self.NodeTemplate.Requirements {
			if req.Name == requirementName {
				isExplicitlyAssigned = true
				break
			}
		}

		if !isExplicitlyAssigned {
			// This is a dangling requirement - defined in node_type but not assigned in node_template
			return reqDef
		}
	}

	return nil
}

//
// RequirementMappings
//

type RequirementMappings map[string]*RequirementMapping
