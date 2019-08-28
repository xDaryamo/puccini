package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// AttributeMapping
//
// Attaches to NotificationDefinition
//

type AttributeMapping struct {
	*Entity `name:"requirement mapping"`

	NodeTemplateName *string `require:"0"`
	AttributeName    *string `require:"1"`
}

func NewAttributeMapping(context *tosca.Context) *AttributeMapping {
	return &AttributeMapping{Entity: NewEntity(context)}
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

//
// AttributeMappings
//

type AttributeMappings []*AttributeMapping
