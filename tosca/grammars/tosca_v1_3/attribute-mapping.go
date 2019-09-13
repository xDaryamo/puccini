package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// AttributeMapping
//
// Attaches to NotificationDefinition
//

type AttributeMapping struct {
	*Entity `name:"attribute mapping"`
	Name    string

	NodeTemplateName *string `require:"0"`
	AttributeName    *string `require:"1"`
}

func NewAttributeMapping(context *tosca.Context) *AttributeMapping {
	return &AttributeMapping{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadAttributeMapping(context *tosca.Context) interface{} {
	self := NewAttributeMapping(context)

	if strings := context.ReadStringListFixed(2); strings != nil {
		self.NodeTemplateName = &(*strings)[0]
		self.AttributeName = &(*strings)[1]
	}

	return self
}

// tosca.Mappable interface
func (self *AttributeMapping) GetKey() string {
	return self.Name
}

//
// AttributeMappings
//

type AttributeMappings map[string]*AttributeMapping

func (self AttributeMappings) Inherit(parent AttributeMappings) {
	for name, attributeMapping := range parent {
		if _, ok := self[name]; !ok {
			self[name] = attributeMapping
		}
	}
}

func (self AttributeMappings) Normalize(n *normal.NodeTemplate, m normal.AttributeMappings) {
	for name, attributeMapping := range self {
		nodeTemplateName := *attributeMapping.NodeTemplateName

		if nodeTemplateName == "SELF" {
			m[name] = n.NewAttributeMapping(*attributeMapping.AttributeName)
		} else {
			if nn, ok := n.ServiceTemplate.NodeTemplates[nodeTemplateName]; ok {
				m[name] = nn.NewAttributeMapping(*attributeMapping.AttributeName)
			}
		}
	}
}
