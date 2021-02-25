package normal

import (
	"encoding/json"
	"math"

	"github.com/fxamacker/cbor/v2"
)

//
// Capability
//

type Capability struct {
	NodeTemplate *NodeTemplate
	Name         string

	Description          string
	Types                Types
	Properties           Constrainables
	Attributes           Constrainables
	MinRelationshipCount uint64
	MaxRelationshipCount uint64
	Location             *Location
}

func (self *NodeTemplate) NewCapability(name string, location *Location) *Capability {
	capability := &Capability{
		NodeTemplate:         self,
		Name:                 name,
		Types:                make(Types),
		Properties:           make(Constrainables),
		Attributes:           make(Constrainables),
		MaxRelationshipCount: math.MaxUint64,
		Location:             location,
	}
	self.Capabilities[name] = capability
	return capability
}

type MarshalableCapability struct {
	Description          string         `json:"description" yaml:"description"`
	Types                Types          `json:"types" yaml:"types"`
	Properties           Constrainables `json:"properties" yaml:"properties"`
	Attributes           Constrainables `json:"attributes" yaml:"attributes"`
	MinRelationshipCount int64          `json:"minRelationshipCount" yaml:"minRelationshipCount"`
	MaxRelationshipCount int64          `json:"maxRelationshipCount" yaml:"maxRelationshipCount"`
	Path                 string         `json:"path" yaml:"path"`
	Location             *Location      `json:"location" yaml:"location"`
}

func (self *Capability) Marshalable() interface{} {
	var minRelationshipCount int64
	var maxRelationshipCount int64
	minRelationshipCount = int64(self.MinRelationshipCount)
	if self.MaxRelationshipCount == math.MaxUint64 {
		maxRelationshipCount = -1
	} else {
		maxRelationshipCount = int64(self.MaxRelationshipCount)
	}

	return &MarshalableCapability{
		Description:          self.Description,
		Types:                self.Types,
		Properties:           self.Properties,
		Attributes:           self.Attributes,
		MinRelationshipCount: minRelationshipCount,
		MaxRelationshipCount: maxRelationshipCount,
		Location:             self.Location,
	}
}

// json.Marshaler interface
func (self *Capability) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Marshalable())
}

// yaml.Marshaler interface
func (self *Capability) MarshalYAML() (interface{}, error) {
	return self.Marshalable(), nil
}

// cbor.Marshaler interface
func (self *Capability) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(self.Marshalable())
}

//
// Capabilities
//

type Capabilities map[string]*Capability
