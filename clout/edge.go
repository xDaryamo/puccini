package clout

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"github.com/tliron/kutil/ard"
)

//
// Edge
//

type Edge struct {
	Metadata   ard.StringMap
	Properties ard.StringMap
	Source     *Vertex
	TargetID   string
	Target     *Vertex
}

func (self *Vertex) NewEdgeTo(target *Vertex) *Edge {
	edge := &Edge{
		Metadata:   make(ard.StringMap),
		Properties: make(ard.StringMap),
		Source:     self,
		Target:     target,
	}
	self.EdgesOut = append(self.EdgesOut, edge)
	target.EdgesIn = append(target.EdgesIn, edge)
	return edge
}

func (self *Vertex) NewEdgeToID(targetId string) *Edge {
	edge := &Edge{
		Metadata:   make(ard.StringMap),
		Properties: make(ard.StringMap),
		Source:     self,
		TargetID:   targetId,
	}
	self.EdgesOut = append(self.EdgesOut, edge)
	return edge
}

func (self *Edge) Remove() {
	self.Source.EdgesOut = self.Source.EdgesOut.Remove(self)
	self.Target.EdgesIn = self.Target.EdgesIn.Remove(self)
}

// Entity interface
func (self *Edge) GetMetadata() ard.StringMap {
	return self.Metadata
}

// Entity interface
func (self *Edge) GetProperties() ard.StringMap {
	return self.Properties
}

type MarshalableEdge struct {
	Metadata   ard.StringMap `yaml:"metadata" cbor:"metadata"`
	Properties ard.StringMap `yaml:"properties" cbor:"properties"`
	TargetID   string        `yaml:"targetID" cbor:"targetID"`
}

type MarshalableEdgeStringMaps struct {
	Metadata   ard.StringMap `json:"metadata"`
	Properties ard.StringMap `json:"properties"`
	TargetID   string        `json:"targetID"`
}

func (self *Edge) Marshalable(stringMaps bool) interface{} {
	var targetID string
	if self.Target != nil {
		targetID = self.Target.ID
	}

	if stringMaps {
		return &MarshalableEdgeStringMaps{
			Metadata:   ard.EnsureStringMaps(self.Metadata),
			Properties: ard.EnsureStringMaps(self.Properties),
			TargetID:   targetID,
		}
	} else {
		return &MarshalableEdge{
			Metadata:   self.Metadata,
			Properties: self.Properties,
			TargetID:   targetID,
		}
	}
}

func (self *Edge) Unmarshal(f func(m *MarshalableEdge) error) error {
	var m MarshalableEdge
	if err := f(&m); err != nil {
		return err
	}
	self.Metadata = m.Metadata
	self.Properties = m.Properties
	self.TargetID = m.TargetID
	return nil
}

// json.Marshaler interface
func (self *Edge) MarshalJSON() ([]byte, error) {
	// JavaScript requires keys to be strings, so we would lose complex keys
	return json.Marshal(self.Marshalable(true))
}

// yaml.Marshaler interface
func (self *Edge) MarshalYAML() (interface{}, error) {
	return self.Marshalable(false), nil
}

// cbor.Marshaler interface
func (self *Edge) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(self.Marshalable(false))
}

// json.Unmarshaler interface
func (self *Edge) UnmarshalJSON(data []byte) error {
	return self.Unmarshal(func(m *MarshalableEdge) error {
		return json.Unmarshal(data, m)
	})
}

// yaml.Unmarshaler interface
func (self *Edge) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return self.Unmarshal(func(m *MarshalableEdge) error {
		return unmarshal(m)
	})
}

// cbor.Unmarshaler interface
func (self *Edge) UnmarshalCBOR(data []byte) error {
	return self.Unmarshal(func(m *MarshalableEdge) error {
		return cbor.Unmarshal(data, m)
	})
}

//
// Edges
//

type Edges []*Edge

func (self Edges) Remove(edge *Edge) Edges {
	for index, e := range self {
		if e == edge {
			return append(self[:index], self[index+1:]...)
		}
	}
	return self
}
