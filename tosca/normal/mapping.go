package normal

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/util"
)

//
// Mapping
//

type Mapping struct {
	NodeTemplate *NodeTemplate
	Relationship *Relationship
	TargetType   string
	Target       string
	Value        Value
}

func NewMapping(targetType string, target string) *Mapping {
	return &Mapping{
		TargetType: targetType,
		Target:     target,
	}
}

func NewMappingValue(targetType string, value Value) *Mapping {
	return &Mapping{
		TargetType: targetType,
		Value:      value,
	}
}

func (self *NodeTemplate) NewMapping(targetType string, target string) *Mapping {
	return &Mapping{
		NodeTemplate: self,
		TargetType:   targetType,
		Target:       target,
	}
}

func (self *Relationship) NewMapping(targetType string, target string) *Mapping {
	return &Mapping{
		Relationship: self,
		TargetType:   targetType,
		Target:       target,
	}
}

type MarshalableMapping struct {
	NodeTemplateName string `json:"nodeTemplateName,omitempty" yaml:"nodeTemplateName,omitempty"`
	TargetType       string `json:"targetType" yaml:"targetType"`
	Target           string `json:"target,omitempty" yaml:"target,omitempty"`
	Value            Value  `json:"value,omitempty" yaml:"value,omitempty"`
}

func (self *Mapping) Marshalable() any {
	if self.NodeTemplate != nil {
		return &MarshalableMapping{
			NodeTemplateName: self.NodeTemplate.Name,
			TargetType:       self.TargetType,
			Target:           self.Target,
			Value:            self.Value,
		}
	} else {
		return &MarshalableMapping{
			TargetType: self.TargetType,
			Target:     self.Target,
			Value:      self.Value,
		}
	}
}

// json.Marshaler interface
func (self *Mapping) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Marshalable())
}

// yaml.Marshaler interface
func (self *Mapping) MarshalYAML() (any, error) {
	return self.Marshalable(), nil
}

// cbor.Marshaler interface
func (self *Mapping) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(self.Marshalable())
}

// msgpack.Marshaler interface
func (self *Mapping) MarshalMsgpack() ([]byte, error) {
	return util.MarshalMessagePack(self.Marshalable())
}

// ard.ToARD interface
func (self *Mapping) ToARD(reflector *ard.Reflector) (any, error) {
	return reflector.Unpack(self.Marshalable())
}

//
// Mappings
//

type Mappings map[string]*Mapping
