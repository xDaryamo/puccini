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

func (self *AttributeMapping) Marshalable() interface{} {
	return &struct {
		NodeTemplateName string `json:"nodeTemplateName" yaml:"nodeTemplateName"`
		AttributeName    string `json:"attributeName" yaml:"attributeName"`
	}{
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
