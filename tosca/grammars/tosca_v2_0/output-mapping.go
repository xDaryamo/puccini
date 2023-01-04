package tosca_v2_0

import (
	"reflect"

	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// OutputMapping
//
// Attaches to notifications and operations
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.15
//

type OutputMapping struct {
	*Entity `name:"output mapping"`
	Name    string

	// The entity name can be a node template name *or* "SELF"
	// If it's "SELF" it could be a node template reference *or* a relationship
	// (but not a group, because a group doesn't have attributes)

	EntityName     *string
	CapabilityName *string
	AttributePath  []string

	NodeTemplate *NodeTemplate           `traverse:"ignore" json:"-" yaml:"-"`
	Relationship *RelationshipAssignment `traverse:"ignore" json:"-" yaml:"-"`
}

func NewOutputMapping(context *tosca.Context) *OutputMapping {
	return &OutputMapping{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadOutputMapping(context *tosca.Context) tosca.EntityPtr {
	self := NewOutputMapping(context)

	if strings := context.ReadStringListMinLength(2); strings != nil {
		self.EntityName = &(*strings)[0]
		self.AttributePath = (*strings)[1:]
	}

	return self
}

// tosca.Mappable interface
func (self *OutputMapping) GetKey() string {
	return self.Name
}

func (self *OutputMapping) RenderForNodeTemplate(nodeTemplate *NodeTemplate) {
	logRender.Debugf("output mapping: %s", self.Name)

	if (self.EntityName == nil) || (len(self.AttributePath) == 0) {
		return
	}

	entityName := *self.EntityName
	switch entityName {
	case "SELF":
		self.setNodeTemplate(nodeTemplate)
	case "SOURCE", "TARGET":
		self.Context.ListChild(0, entityName).ReportValueInvalid("modelable entity name", "cannot be used in node templates")
	default:
		self.setNodeTemplateByName(entityName)
	}
}

func (self *OutputMapping) RenderForRelationship(relationship *RelationshipAssignment) {
	logRender.Debugf("output mapping: %s", self.Name)

	if (self.EntityName == nil) || (len(self.AttributePath) == 0) {
		return
	}

	switch entityName := *self.EntityName; entityName {
	case "SELF", "SOURCE", "TARGET":
		self.setRelationship(relationship)
	default:
		self.setNodeTemplateByName(entityName)
	}
}

func (self *OutputMapping) RenderForGroup() {
	logRender.Debugf("output mapping: %s", self.Name)

	if (self.EntityName == nil) || (len(self.AttributePath) == 0) {
		return
	}

	switch entityName := *self.EntityName; entityName {
	case "SELF", "SOURCE", "TARGET":
		self.Context.ListChild(0, entityName).ReportValueInvalid("modelable entity name", "cannot be used in groups")
	default:
		self.setNodeTemplateByName(entityName)
	}
}

func (self *OutputMapping) setNodeTemplate(nodeTemplate *NodeTemplate) {
	self.NodeTemplate = nodeTemplate

	// Attributes should already have been rendered
	attributeName := self.AttributePath[0]
	if _, ok := self.NodeTemplate.Attributes[attributeName]; !ok {
		self.Context.ListChild(1, attributeName).ReportReferenceNotFound("attribute", self.NodeTemplate)
	}
}

func (self *OutputMapping) setNodeTemplateByName(nodeTemplateName string) {
	var nodeTemplateType *NodeTemplate
	if nodeTemplate, ok := self.Context.Namespace.LookupForType(nodeTemplateName, reflect.TypeOf(nodeTemplateType)); ok {
		self.setNodeTemplate(nodeTemplate.(*NodeTemplate))
	} else {
		self.Context.ListChild(0, nodeTemplateName).ReportUnknown("node template")
		return
	}
}

func (self *OutputMapping) setRelationship(relationship *RelationshipAssignment) {
	self.Relationship = relationship

	// Attributes should already have been rendered
	attributeName := self.AttributePath[0]
	if _, ok := self.Relationship.Attributes[attributeName]; !ok {
		self.Context.ListChild(1, attributeName).ReportReferenceNotFound("attribute", self.Relationship)
	}
}

func (self *OutputMapping) Normalize(normalOutputs normal.Values) {
	if (self.EntityName == nil) || (len(self.AttributePath) == 0) {
		return
	}

	list := normal.NewList(len(self.AttributePath) + 1)
	meta := normal.NewValueMeta()
	meta.Type = "list"
	meta.Element = normal.NewValueMeta()
	meta.Element.Type = "string"
	list.SetMeta(meta)

	switch entityName := *self.EntityName; entityName {
	case "SOURCE", "TARGET":
		// These can only be retrieved via a function call
		row, column := self.Context.GetLocation()
		arguments := []any{normal.NewPrimitive(entityName)}
		list.Set(0, normal.NewFunctionCall(tosca.NewFunctionCall("tosca.function._get_modelable_entity", arguments, self.Context.URL.String(), row, column, self.Context.Path.String())))
	default:
		list.Set(0, normal.NewPrimitive(entityName))
	}

	for index, pathElement := range self.AttributePath {
		list.Set(index+1, normal.NewPrimitive(pathElement))
	}

	normalOutputs[self.Name] = list
}

//
// OutputMappings
//

type OutputMappings map[string]*OutputMapping

func (self OutputMappings) CopyUnassigned(outputMappings OutputMappings) {
	for key, outputMapping := range outputMappings {
		if _, ok := self[key]; !ok {
			self[key] = outputMapping
		}
	}
}

func (self OutputMappings) Inherit(parent OutputMappings) {
	for name, outputMapping := range parent {
		if _, ok := self[name]; !ok {
			self[name] = outputMapping
		}
	}
}

func (self OutputMappings) RenderForNodeTemplate(nodeTemplate *NodeTemplate) {
	for _, outputMapping := range self {
		outputMapping.RenderForNodeTemplate(nodeTemplate)
	}
}

func (self OutputMappings) RenderForRelationship(relationship *RelationshipAssignment) {
	for _, outputMapping := range self {
		outputMapping.RenderForRelationship(relationship)
	}
}

func (self OutputMappings) RenderForGroup() {
	for _, outputMapping := range self {
		outputMapping.RenderForGroup()
	}
}

func (self OutputMappings) Normalize(normalOutputs normal.Values) {
	for _, outputMapping := range self {
		outputMapping.Normalize(normalOutputs)
	}
}
