package tosca_v2_0

import (
	"strconv"
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// PropertyMapping
//
// Attaches to SubstitutionMappings
//
// [TOSCA-v2.0] @ 15.2
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.8
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.8
//

type PropertyMapping struct {
	*Entity `name:"property mapping"`
	Name    string

	InputName *string

	// TOSCA 2.0 - Direct property mapping
	PropertyName *string

	// TOSCA 2.0 - Capability property mapping
	CapabilityName         *string
	CapabilityPropertyName *string

	// TOSCA 2.0 - Relationship property mapping
	RequirementName          *string
	RelationshipIndex        *int // nil means 0, -1 means ALL
	RelationshipPropertyName *string

	// Deprecated TOSCA 1.x fields
	NodeTemplateName *string // deprecated in TOSCA 1.3

	InputDefinition *ParameterDefinition `traverse:"ignore" json:"-" yaml:"-"`
	NodeTemplate    *NodeTemplate        `traverse:"ignore" json:"-" yaml:"-"` // deprecated in TOSCA 1.3
	Value           *Value               `traverse:"ignore" json:"-" yaml:"-"` // deprecated in TOSCA 1.3
}

func NewPropertyMapping(context *parsing.Context) *PropertyMapping {
	return &PropertyMapping{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadPropertyMapping(context *parsing.Context) parsing.EntityPtr {
	self := NewPropertyMapping(context)

	if self.tryParseMultilineMapping(context) {
		return self
	}

	if self.tryParseListMapping(context) {
		return self
	}

	if self.tryParseStringMapping(context) {
		return self
	}

	// Fallback to constant value (deprecated in TOSCA 1.3)
	self.Value = ReadValue(context).(*Value)
	return self
}

func (self *PropertyMapping) tryParseMultilineMapping(context *parsing.Context) bool {
	if !self.isMultilineMapping(context) {
		return false
	}

	parsedList := self.parseMultilineKey(context.Name)
	return self.parseStructuredMapping(context, parsedList)
}

func (self *PropertyMapping) isMultilineMapping(context *parsing.Context) bool {
	return strings.Contains(context.Name, "¶") &&
		(strings.Contains(context.Name, "CAPABILITY") || strings.Contains(context.Name, "RELATIONSHIP"))
}

func (self *PropertyMapping) parseMultilineKey(name string) []string {
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

func (self *PropertyMapping) parseStructuredMapping(context *parsing.Context, parsedList []string) bool {
	if len(parsedList) < 3 {
		return false
	}

	inputName := context.ReadString()
	if inputName == nil {
		return false
	}

	switch parsedList[0] {
	case "CAPABILITY":
		return self.parseCapabilityMapping(parsedList, inputName)
	case "RELATIONSHIP":
		return self.parseRelationshipMapping(context, parsedList, inputName)
	default:
		return false
	}
}

func (self *PropertyMapping) parseCapabilityMapping(parsedList []string, inputName *string) bool {
	if len(parsedList) != 3 {
		return false
	}

	self.CapabilityName = &parsedList[1]
	self.CapabilityPropertyName = &parsedList[2]
	self.InputName = inputName
	return true
}

func (self *PropertyMapping) parseRelationshipMapping(context *parsing.Context, parsedList []string, inputName *string) bool {
	switch len(parsedList) {
	case 3:
		return self.parseRelationshipMappingWithDefaultIndex(parsedList, inputName)
	case 4:
		return self.parseRelationshipMappingWithIndex(context, parsedList, inputName)
	default:
		return false
	}
}

func (self *PropertyMapping) parseRelationshipMappingWithDefaultIndex(parsedList []string, inputName *string) bool {
	self.RequirementName = &parsedList[1]
	idx := 0
	self.RelationshipIndex = &idx
	self.RelationshipPropertyName = &parsedList[2]
	self.InputName = inputName
	return true
}

func (self *PropertyMapping) parseRelationshipMappingWithIndex(context *parsing.Context, parsedList []string, inputName *string) bool {
	self.RequirementName = &parsedList[1]

	if !self.parseRelationshipIndex(context, parsedList[2]) {
		return false
	}

	self.RelationshipPropertyName = &parsedList[3]
	self.InputName = inputName
	return true
}

func (self *PropertyMapping) parseRelationshipIndex(context *parsing.Context, indexStr string) bool {
	if indexStr == "ALL" {
		idx := -1
		self.RelationshipIndex = &idx
		return true
	}

	if idx, err := strconv.Atoi(indexStr); err == nil {
		self.RelationshipIndex = &idx
		return true
	}

	context.ReportValueMalformed("property mapping", "relationship index must be a number or 'ALL'")
	return false
}

func (self *PropertyMapping) tryParseListMapping(context *parsing.Context) bool {
	if !context.Is(ard.TypeList) {
		return false
	}

	strings := context.ReadStringList()
	if strings == nil {
		return false
	}

	switch len(*strings) {
	case 1:
		return self.parseSimpleInputMapping(*strings)
	case 2:
		return self.parseTwoParameterMapping(context, *strings)
	case 3:
		return self.parseThreeParameterMapping(context, *strings)
	case 4:
		return self.parseFourParameterMapping(context, *strings)
	default:
		context.ReportValueMalformed("property mapping", "property mapping must be 1-4 parameters")
		return false
	}
}

func (self *PropertyMapping) parseSimpleInputMapping(strings []string) bool {
	self.InputName = &strings[0]
	return true
}

func (self *PropertyMapping) parseTwoParameterMapping(context *parsing.Context, strings []string) bool {
	if strings[0] == "CAPABILITY" || strings[0] == "RELATIONSHIP" {
		context.ReportValueMalformed("property mapping", "CAPABILITY and RELATIONSHIP mappings require more parameters")
		return false
	}

	// Deprecated TOSCA 1.x format: [<node_template_name>, <property_name>]
	self.NodeTemplateName = &strings[0]
	self.PropertyName = &strings[1]
	return true
}

func (self *PropertyMapping) parseThreeParameterMapping(context *parsing.Context, strings []string) bool {
	if strings[0] != "CAPABILITY" {
		context.ReportValueMalformed("property mapping", "invalid 3-parameter format")
		return false
	}

	self.CapabilityName = &strings[1]
	self.CapabilityPropertyName = &strings[2]

	if inputName := context.ReadString(); inputName != nil {
		self.InputName = inputName
	}

	return true
}

func (self *PropertyMapping) parseFourParameterMapping(context *parsing.Context, strings []string) bool {
	if strings[0] != "RELATIONSHIP" {
		context.ReportValueMalformed("property mapping", "invalid 4-parameter format")
		return false
	}

	self.RequirementName = &strings[1]

	if !self.parseRelationshipIndex(context, strings[2]) {
		return false
	}

	self.RelationshipPropertyName = &strings[3]

	if inputName := context.ReadString(); inputName != nil {
		self.InputName = inputName
	}

	return true
}

func (self *PropertyMapping) tryParseStringMapping(context *parsing.Context) bool {
	if !context.Is(ard.TypeString) {
		return false
	}

	inputName := context.ReadString()
	if inputName == nil {
		return false
	}

	self.InputName = inputName
	self.PropertyName = &self.Name
	return true
}

// ([parsing.Mappable] interface)
func (self *PropertyMapping) GetKey() string {
	return self.Name
}

func (self *PropertyMapping) Render(inputDefinitions ParameterDefinitions) {
	logRender.Debug("property mapping")

	if self.InputName != nil {
		self.renderInputMapping(inputDefinitions)
	} else if self.isDeprecatedNodeTemplateMapping() {
		self.renderDeprecatedNodeTemplateMapping()
	}
}

func (self *PropertyMapping) renderInputMapping(inputDefinitions ParameterDefinitions) {
	inputName := *self.InputName

	var ok bool
	if self.InputDefinition, ok = inputDefinitions[inputName]; !ok {
		self.Context.ListChild(0, inputName).ReportUnknown("input")
	}
}

func (self *PropertyMapping) isDeprecatedNodeTemplateMapping() bool {
	return self.NodeTemplateName != nil && self.PropertyName != nil
}

func (self *PropertyMapping) renderDeprecatedNodeTemplateMapping() {
	nodeTemplateName := *self.NodeTemplateName

	nodeTemplate, ok := self.Context.Namespace.LookupForType(nodeTemplateName, nodeTemplatePtrType)
	if !ok {
		self.Context.ListChild(0, nodeTemplateName).ReportUnknown("node template")
		return
	}

	self.NodeTemplate = nodeTemplate.(*NodeTemplate)
	self.NodeTemplate.Render()

	self.renderNodeTemplateProperty()
}

func (self *PropertyMapping) renderNodeTemplateProperty() {
	propertyName := *self.PropertyName

	var ok bool
	if self.Value, ok = self.NodeTemplate.Properties[propertyName]; !ok {
		self.Context.ListChild(1, propertyName).ReportReferenceNotFound("property", self.NodeTemplate)
	}
}

// Utility methods for type checking
func (self *PropertyMapping) IsCapabilityMapping() bool {
	return self.CapabilityName != nil && self.CapabilityPropertyName != nil
}

func (self *PropertyMapping) IsRelationshipMapping() bool {
	return self.RequirementName != nil && self.RelationshipPropertyName != nil
}

func (self *PropertyMapping) IsDirectPropertyMapping() bool {
	return self.PropertyName != nil && !self.IsCapabilityMapping() && !self.IsRelationshipMapping()
}

func (self *PropertyMapping) IsInputMapping() bool {
	return self.InputName != nil
}

func (self *PropertyMapping) IsValueMapping() bool {
	return self.Value != nil
}

func (self *PropertyMapping) IsAllRelationshipMapping() bool {
	return self.IsRelationshipMapping() && self.RelationshipIndex != nil && *self.RelationshipIndex == -1
}

func (self *PropertyMapping) GetRelationshipIndexValue() int {
	if self.RelationshipIndex == nil {
		return 0
	}
	return *self.RelationshipIndex
}

//
// PropertyMappings
//

type PropertyMappings map[string]*PropertyMapping

func (self PropertyMappings) Render(inputDefinitions ParameterDefinitions) {
	for _, mapping := range self {
		mapping.Render(inputDefinitions)
	}
}

func (self PropertyMappings) GetCapabilityMappings() PropertyMappings {
	result := make(PropertyMappings)
	for key, mapping := range self {
		if mapping.IsCapabilityMapping() {
			result[key] = mapping
		}
	}
	return result
}

func (self PropertyMappings) GetRelationshipMappings() PropertyMappings {
	result := make(PropertyMappings)
	for key, mapping := range self {
		if mapping.IsRelationshipMapping() {
			result[key] = mapping
		}
	}
	return result
}

func (self PropertyMappings) GetDirectPropertyMappings() PropertyMappings {
	result := make(PropertyMappings)
	for key, mapping := range self {
		if mapping.IsDirectPropertyMapping() {
			result[key] = mapping
		}
	}
	return result
}
