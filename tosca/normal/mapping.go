package normal

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
)

//
// Mapping
//

type Mapping struct {
	NodeTemplate *NodeTemplate
	Name         string
}

func (self *NodeTemplate) NewMapping(name string) *Mapping {
	return &Mapping{
		NodeTemplate: self,
		Name:         name,
	}
}

type MarshalableAttributeMapping struct {
	NodeTemplateName string `json:"nodeTemplateName" yaml:"nodeTemplateName"`
	Name             string `json:"name" yaml:"name"`
}

func (self *Mapping) Marshalable() interface{} {
	return &MarshalableAttributeMapping{
		NodeTemplateName: self.NodeTemplate.Name,
		Name:             self.Name,
	}
}

// json.Marshaler interface
func (self *Mapping) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Marshalable())
}

// yaml.Marshaler interface
func (self *Mapping) MarshalYAML() (interface{}, error) {
	return self.Marshalable(), nil
}

// cbor.Marshaler interface
func (self *Mapping) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(self.Marshalable())
}

//
// Mappings
//

type Mappings map[*NodeTemplate]*Mapping

//
// Outputs
//

type Outputs map[string]*Mapping
