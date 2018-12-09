package clout

import (
	"encoding/json"

	"github.com/tliron/puccini/ard"
)

//
// Edge
//

type Edge struct {
	Metadata   ard.Map
	Properties ard.Map
	Source     *Vertex
	TargetID   string
	Target     *Vertex
}

func (self *Vertex) NewEdgeTo(target *Vertex) *Edge {
	edge := &Edge{
		Metadata:   make(ard.Map),
		Properties: make(ard.Map),
		Source:     self,
		Target:     target,
	}
	self.EdgesOut = append(self.EdgesOut, edge)
	target.EdgesIn = append(target.EdgesIn, edge)
	return edge
}

// Entity interface
func (self *Edge) GetMetadata() ard.Map {
	return self.Metadata
}

// Entity interface
func (self *Edge) GetProperties() ard.Map {
	return self.Properties
}

type MarshalableEdge struct {
	Metadata   ard.Map `json:"metadata" yaml:"metadata"`
	Properties ard.Map `json:"properties" yaml:"properties"`
	TargetID   string  `json:"targetID" yaml:"targetID"`
}

func (self *Edge) Marshalable() interface{} {
	var targetID string
	if self.Target != nil {
		targetID = self.Target.ID
	}

	return &MarshalableEdge{
		Metadata:   self.Metadata,
		Properties: self.Properties,
		TargetID:   targetID,
	}
}

func (self *Edge) Unmarshal(f func(m *MarshalableEdge) error) error {
	var m MarshalableEdge
	err := f(&m)
	if err != nil {
		return err
	}
	self.Metadata = m.Metadata
	self.Properties = m.Properties
	self.TargetID = m.TargetID
	return nil
}

// json.Marshaler interface
func (self *Edge) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Marshalable())
}

// yaml.Marshaler interface
func (self *Edge) MarshalYAML() (interface{}, error) {
	return self.Marshalable(), nil
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

//
// Edges
//

type Edges []*Edge
