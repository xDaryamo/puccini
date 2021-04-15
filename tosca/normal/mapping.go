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
	TargetType   string
	Target       string
}

func NewMapping(targetType string, target string) *Mapping {
	return &Mapping{
		TargetType: targetType,
		Target:     target,
	}
}

func (self *NodeTemplate) NewMapping(targetType string, target string) *Mapping {
	return &Mapping{
		NodeTemplate: self,
		TargetType:   targetType,
		Target:       target,
	}
}

type MarshalableMapping struct {
	NodeTemplateName string `json:"nodeTemplateName,omitempty" yaml:"nodeTemplateName,omitempty"`
	TargetType       string `json:"targetType" yaml:"targetType"`
	Target           string `json:"target" yaml:"target"`
}

func (self *Mapping) Marshalable() interface{} {
	if self.NodeTemplate != nil {
		return &MarshalableMapping{
			NodeTemplateName: self.NodeTemplate.Name,
			TargetType:       self.TargetType,
			Target:           self.Target,
		}
	} else {
		return &MarshalableMapping{
			TargetType: self.TargetType,
			Target:     self.Target,
		}
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

type Mappings map[string]*Mapping
