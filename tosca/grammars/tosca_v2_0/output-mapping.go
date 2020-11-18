package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// OutputMapping
//
// Attaches to notifications and operations
//

type OutputMapping struct {
	*Entity `name:"output mapping"`
	Name    string

	NodeTemplateName *string `require:"0"`
	AttributeName    *string `require:"1"`
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

	if strings := context.ReadStringListFixed(2); strings != nil {
		self.NodeTemplateName = &(*strings)[0]
		self.AttributeName = &(*strings)[1]
	}

	return self
}

// tosca.Mappable interface
func (self *OutputMapping) GetKey() string {
	return self.Name
}

func (self *OutputMapping) Normalize(name string, normalNodeTemplate *normal.NodeTemplate, normalAttributeMappings normal.AttributeMappings) {
	if (self.NodeTemplateName == nil) || (self.AttributeName == nil) {
		return
	}

	nodeTemplateName := *self.NodeTemplateName
	if nodeTemplateName == "SELF" {
		normalAttributeMappings[name] = normalNodeTemplate.NewAttributeMapping(*self.AttributeName)
	} else {
		if normalOutputNodeTemplate, ok := normalNodeTemplate.ServiceTemplate.NodeTemplates[nodeTemplateName]; ok {
			normalAttributeMappings[name] = normalOutputNodeTemplate.NewAttributeMapping(*self.AttributeName)
		} else {
			self.Context.ListChild(0, nodeTemplateName).ReportUnknown("node template")
		}
	}
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

func (self OutputMappings) Normalize(normalNodeTemplate *normal.NodeTemplate, normalAttributeMappings normal.AttributeMappings) {
	for name, outputMapping := range self {
		outputMapping.Normalize(name, normalNodeTemplate, normalAttributeMappings)
	}
}
