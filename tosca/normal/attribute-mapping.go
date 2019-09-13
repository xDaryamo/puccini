package normal

import (
	"encoding/json"
)

//
// AttributeMapping
//

type AttributeMapping struct {
	NodeTemplate  *NodeTemplate
	AttributeName string
}

func (self *NodeTemplate) NewAttributeMapping(attributeName string) *AttributeMapping {
	return &AttributeMapping{
		NodeTemplate:  self,
		AttributeName: attributeName,
	}
}

type MarshalableAttributeMapping struct {
	NodeTemplateName string `json:"nodeTemplateName" yaml:"nodeTemplateName"`
	AttributeName    string `json:"attributeName" yaml:"attributeName"`
}

func (self *AttributeMapping) Marshalable() interface{} {
	return &MarshalableAttributeMapping{
		NodeTemplateName: self.NodeTemplate.Name,
		AttributeName:    self.AttributeName,
	}
}

// json.Marshaler interface
func (self *AttributeMapping) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Marshalable())
}

// yaml.Marshaler interface
func (self *AttributeMapping) MarshalYAML() (interface{}, error) {
	return self.Marshalable(), nil
}

//
// AttributeMappings
//

type AttributeMappings map[string]*AttributeMapping
