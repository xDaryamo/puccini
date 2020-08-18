package tosca_v2_0

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
func ReadNotificationOutput(context *tosca.Context) tosca.EntityPtr {
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

func (self NotificationOutputs) CopyUnassigned(outputs NotificationOutputs) {
	for key, output := range outputs {
		if _, ok := self[key]; !ok {
			self[key] = output
		}
	}
}

func (self NotificationOutputs) Inherit(parent NotificationOutputs) {
	for name, notificationOutput := range parent {
		if _, ok := self[name]; !ok {
			self[name] = notificationOutput
		}
	}
}

func (self NotificationOutputs) Normalize(normalNodeTemplate *normal.NodeTemplate, normalAttributeMappings normal.AttributeMappings) {
	for name, notificationOutput := range self {
		nodeTemplateName := *notificationOutput.NodeTemplateName

		if nodeTemplateName == "SELF" {
			normalAttributeMappings[name] = normalNodeTemplate.NewAttributeMapping(*notificationOutput.AttributeName)
		} else {
			if normalOutputNodeTemplate, ok := normalNodeTemplate.ServiceTemplate.NodeTemplates[nodeTemplateName]; ok {
				normalAttributeMappings[name] = normalOutputNodeTemplate.NewAttributeMapping(*notificationOutput.AttributeName)
			}
		}
	}
}
