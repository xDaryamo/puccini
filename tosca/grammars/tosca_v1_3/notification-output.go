package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// NotificationOutput
//
// Attaches to NotificationDefinition
//

type NotificationOutput struct {
	*Entity `name:"notification output"`
	Name    string

	NodeTemplateName *string `require:"0"`
	AttributeName    *string `require:"1"`
}

func NewNotificationOutput(context *tosca.Context) *NotificationOutput {
	return &NotificationOutput{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadNotificationOutput(context *tosca.Context) interface{} {
	self := NewNotificationOutput(context)

	if strings := context.ReadStringListFixed(2); strings != nil {
		self.NodeTemplateName = &(*strings)[0]
		self.AttributeName = &(*strings)[1]
	}

	return self
}

// tosca.Mappable interface
func (self *NotificationOutput) GetKey() string {
	return self.Name
}

//
// NotificationOutputs
//

type NotificationOutputs map[string]*NotificationOutput

func (self NotificationOutputs) Inherit(parent NotificationOutputs) {
	for name, notificationOutput := range parent {
		if _, ok := self[name]; !ok {
			self[name] = notificationOutput
		}
	}
}

func (self NotificationOutputs) Normalize(n *normal.NodeTemplate, m normal.AttributeMappings) {
	for name, notificationOutput := range self {
		nodeTemplateName := *notificationOutput.NodeTemplateName

		if nodeTemplateName == "SELF" {
			m[name] = n.NewAttributeMapping(*notificationOutput.AttributeName)
		} else {
			if nn, ok := n.ServiceTemplate.NodeTemplates[nodeTemplateName]; ok {
				m[name] = nn.NewAttributeMapping(*notificationOutput.AttributeName)
			}
		}
	}
}
