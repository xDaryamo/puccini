package tosca_v2_0

import (
	"strconv"
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// AttributeMapping
//
// Attaches to SubstitutionMappings
//
// [TOSCA-v2.0] @ 15.3
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.15
//

type AttributeMapping struct {
	*Entity `name:"attribute mapping"`
	Name    string

	OutputName *string

	// TOSCA 2.0 - Direct attribute mapping
	AttributeName *string

	// TOSCA 2.0 - Capability attribute mapping
	CapabilityName        *string
	CapabilityAttributeName *string

	// TOSCA 2.0 - Relationship attribute mapping
	RequirementName         *string
	RelationshipIndex       *int // nil means 0, -1 means ALL
	RelationshipAttributeName *string

	// Deprecated TOSCA 1.x fields
	NodeTemplateName *string // deprecated in TOSCA 1.3

	OutputDefinition *ParameterDefinition `traverse:"ignore" json:"-" yaml:"-"`
	NodeTemplate     *NodeTemplate        `traverse:"ignore" json:"-" yaml:"-"` // deprecated in TOSCA 1.3
	Attribute        *Value               `traverse:"ignore" json:"-" yaml:"-"` // deprecated in TOSCA 1.3
}

func NewAttributeMapping(context *parsing.Context) *AttributeMapping {
	return &AttributeMapping{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadAttributeMapping(context *parsing.Context) parsing.EntityPtr {
	self := NewAttributeMapping(context)

	if self.tryParseMultilineMapping(context) {
		return self
	}

	if self.tryParseListMapping(context) {
		return self
	}

	if self.tryParseStringMapping(context) {
		return self
	}

	// Fallback to deprecated TOSCA 1.x format: [<node_template_name>, <attribute_name>]
	if strings := context.ReadStringListFixed(2); strings != nil {
		self.NodeTemplateName = &(*strings)[0]
		self.AttributeName = &(*strings)[1]
	}

	return self
}

func (self *AttributeMapping) tryParseMultilineMapping(context *parsing.Context) bool {
	if !self.isMultilineMapping(context) {
		return false
	}

	parsedList := self.parseMultilineKey(context.Name)
	return self.parseStructuredMapping(context, parsedList)
}

func (self *AttributeMapping) isMultilineMapping(context *parsing.Context) bool {
	return strings.Contains(context.Name, "¶") &&
		(strings.Contains(context.Name, "CAPABILITY") || strings.Contains(context.Name, "RELATIONSHIP"))
}

func (self *AttributeMapping) parseMultilineKey(name string) []string {
	lines := strings.Split(strings.ReplaceAll(name, "¶", "\n"), "\n")
	var parsedList []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") {
			parsedList = append(parsedList, strings.TrimPrefix(line, "- "))
		}
	}

	return parsedList
}

func (self *AttributeMapping) parseStructuredMapping(context *parsing.Context, parsedList []string) bool {
	if len(parsedList) < 3 {
		return false
	}

	outputName := context.ReadString()
	if outputName == nil {
		return false
	}

	switch parsedList[0] {
	case "CAPABILITY":
		return self.parseCapabilityMapping(parsedList, outputName)
	case "RELATIONSHIP":
		return self.parseRelationshipMapping(context, parsedList, outputName)
	default:
		return false
	}
}

func (self *AttributeMapping) parseCapabilityMapping(parsedList []string, outputName *string) bool {
	if len(parsedList) != 3 {
		return false
	}

	self.CapabilityName = &parsedList[1]
	self.CapabilityAttributeName = &parsedList[2]
	self.OutputName = outputName
	return true
}

func (self *AttributeMapping) parseRelationshipMapping(context *parsing.Context, parsedList []string, outputName *string) bool {
	switch len(parsedList) {
	case 3:
		return self.parseRelationshipMappingWithDefaultIndex(parsedList, outputName)
	case 4:
		return self.parseRelationshipMappingWithIndex(context, parsedList, outputName)
	default:
		return false
	}
}

func (self *AttributeMapping) parseRelationshipMappingWithDefaultIndex(parsedList []string, outputName *string) bool {
	self.RequirementName = &parsedList[1]
	idx := 0
	self.RelationshipIndex = &idx
	self.RelationshipAttributeName = &parsedList[2]
	self.OutputName = outputName
	return true
}

func (self *AttributeMapping) parseRelationshipMappingWithIndex(context *parsing.Context, parsedList []string, outputName *string) bool {
	self.RequirementName = &parsedList[1]

	if !self.parseRelationshipIndex(context, parsedList[2]) {
		return false
	}

	self.RelationshipAttributeName = &parsedList[3]
	self.OutputName = outputName
	return true
}

func (self *AttributeMapping) parseRelationshipIndex(context *parsing.Context, indexStr string) bool {
	if indexStr == "ALL" {
		idx := -1
		self.RelationshipIndex = &idx
		return true
	}

	if idx, err := strconv.Atoi(indexStr); err == nil {
		self.RelationshipIndex = &idx
		return true
	}

	context.ReportValueMalformed("attribute mapping", "relationship index must be a number or 'ALL'")
	return false
}

func (self *AttributeMapping) tryParseListMapping(context *parsing.Context) bool {
	if !context.Is(ard.TypeList) {
		return false
	}

	strings := context.ReadStringList()
	if strings == nil {
		return false
	}

	switch len(*strings) {
	case 1:
		return self.parseSimpleOutputMapping(*strings)
	case 2:
		return self.parseTwoParameterMapping(context, *strings)
	case 3:
		return self.parseThreeParameterMapping(context, *strings)
	case 4:
		return self.parseFourParameterMapping(context, *strings)
	default:
		context.ReportValueMalformed("attribute mapping", "attribute mapping must be 1-4 parameters")
		return false
	}
}

func (self *AttributeMapping) parseSimpleOutputMapping(strings []string) bool {
	self.OutputName = &strings[0]
	return true
}

func (self *AttributeMapping) parseTwoParameterMapping(context *parsing.Context, strings []string) bool {
	if strings[0] == "CAPABILITY" || strings[0] == "RELATIONSHIP" {
		context.ReportValueMalformed("attribute mapping", "CAPABILITY and RELATIONSHIP mappings require more parameters")
		return false
	}

	// Deprecated TOSCA 1.x format: [<node_template_name>, <attribute_name>]
	self.NodeTemplateName = &strings[0]
	self.AttributeName = &strings[1]
	return true
}

func (self *AttributeMapping) parseThreeParameterMapping(context *parsing.Context, strings []string) bool {
	if strings[0] != "CAPABILITY" {
		context.ReportValueMalformed("attribute mapping", "invalid 3-parameter format")
		return false
	}

	self.CapabilityName = &strings[1]
	self.CapabilityAttributeName = &strings[2]

	if outputName := context.ReadString(); outputName != nil {
		self.OutputName = outputName
	}

	return true
}

func (self *AttributeMapping) parseFourParameterMapping(context *parsing.Context, strings []string) bool {
	if strings[0] != "RELATIONSHIP" {
		context.ReportValueMalformed("attribute mapping", "invalid 4-parameter format")
		return false
	}

	self.RequirementName = &strings[1]

	if !self.parseRelationshipIndex(context, strings[2]) {
		return false
	}

	self.RelationshipAttributeName = &strings[3]

	if outputName := context.ReadString(); outputName != nil {
		self.OutputName = outputName
	}

	return true
}

func (self *AttributeMapping) tryParseStringMapping(context *parsing.Context) bool {
	if !context.Is(ard.TypeString) {
		return false
	}

	outputName := context.ReadString()
	if outputName == nil {
		return false
	}

	self.OutputName = outputName
	self.AttributeName = &self.Name
	return true
}

// ([parsing.Mappable] interface)
func (self *AttributeMapping) GetKey() string {
	return self.Name
}

func (self *AttributeMapping) EnsureRender() {
	logRender.Debug("attribute mapping")

	if self.OutputName != nil {
		self.renderOutputMapping()
	} else if self.isDeprecatedNodeTemplateMapping() {
		self.renderDeprecatedNodeTemplateMapping()
	}
}

func (self *AttributeMapping) renderOutputMapping() {
	// For TOSCA 2.0, we don't validate outputs here as they belong to the service template
	// The validation should be done in the substitution mappings context
}

func (self *AttributeMapping) isDeprecatedNodeTemplateMapping() bool {
	return self.NodeTemplateName != nil && self.AttributeName != nil
}

func (self *AttributeMapping) renderDeprecatedNodeTemplateMapping() {
	nodeTemplateName := *self.NodeTemplateName

	nodeTemplate, ok := self.Context.Namespace.LookupForType(nodeTemplateName, nodeTemplatePtrType)
	if !ok {
		self.Context.ListChild(0, nodeTemplateName).ReportUnknown("node template")
		return
	}

	self.NodeTemplate = nodeTemplate.(*NodeTemplate)
	self.NodeTemplate.Render()

	self.renderNodeTemplateAttribute()
}

func (self *AttributeMapping) renderNodeTemplateAttribute() {
	attributeName := *self.AttributeName

	var ok bool
	if self.Attribute, ok = self.NodeTemplate.Attributes[attributeName]; !ok {
		self.Context.ListChild(1, attributeName).ReportReferenceNotFound("attribute", self.NodeTemplate)
	}
}

// Utility methods for type checking
func (self *AttributeMapping) IsCapabilityMapping() bool {
	return self.CapabilityName != nil && self.CapabilityAttributeName != nil
}

func (self *AttributeMapping) IsRelationshipMapping() bool {
	return self.RequirementName != nil && self.RelationshipAttributeName != nil
}

func (self *AttributeMapping) IsDirectAttributeMapping() bool {
	return self.AttributeName != nil && !self.IsCapabilityMapping() && !self.IsRelationshipMapping()
}

func (self *AttributeMapping) IsOutputMapping() bool {
	return self.OutputName != nil
}

func (self *AttributeMapping) IsAllRelationshipMapping() bool {
	return self.IsRelationshipMapping() && self.RelationshipIndex != nil && *self.RelationshipIndex == -1
}

func (self *AttributeMapping) GetRelationshipIndexValue() int {
	if self.RelationshipIndex == nil {
		return 0
	}
	return *self.RelationshipIndex
}

//
// AttributeMappings
//

type AttributeMappings map[string]*AttributeMapping

func (self AttributeMappings) EnsureRender() {
	for _, mapping := range self {
		mapping.EnsureRender()
	}
}

func (self AttributeMappings) GetCapabilityMappings() AttributeMappings {
	result := make(AttributeMappings)
	for key, mapping := range self {
		if mapping.IsCapabilityMapping() {
			result[key] = mapping
		}
	}
	return result
}

func (self AttributeMappings) GetRelationshipMappings() AttributeMappings {
	result := make(AttributeMappings)
	for key, mapping := range self {
		if mapping.IsRelationshipMapping() {
			result[key] = mapping
		}
	}
	return result
}

func (self AttributeMappings) GetDirectAttributeMappings() AttributeMappings {
	result := make(AttributeMappings)
	for key, mapping := range self {
		if mapping.IsDirectAttributeMapping() {
			result[key] = mapping
		}
	}
	return result
}
