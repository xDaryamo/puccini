package normal

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
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

// cbor.Marshaler interface
func (self *AttributeMapping) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(self.Marshalable())
}

//
// AttributeMappings
//

type AttributeMappings map[string]*AttributeMapping
