package clout

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/format"
)

const Version = "1.0"

//
// Clout
//

type Clout struct {
	Version    string   `json:"version" yaml:"version"`
	Metadata   ard.Map  `json:"metadata" yaml:"metadata"`
	Properties ard.Map  `json:"properties" yaml:"properties"`
	Vertexes   Vertexes `json:"vertexes" yaml:"vertexes"`
}

func NewClout() *Clout {
	return &Clout{
		Version:    Version,
		Metadata:   make(ard.Map),
		Properties: make(ard.Map),
		Vertexes:   make(Vertexes),
	}
}

func (self *Clout) Resolve() error {
	if self.Version == "" {
		return fmt.Errorf("no Clout \"Version\"")
	}
	if self.Version != Version {
		return fmt.Errorf("unsupported Clout version: \"%s\"", self.Version)
	}

	self.Metadata = ard.EnsureMap(self.Metadata)
	self.Properties = ard.EnsureMap(self.Properties)

	for key, v := range self.Vertexes {
		v.ID = key
		v.Metadata = ard.EnsureMap(v.Metadata)
		v.Properties = ard.EnsureMap(v.Properties)

		for _, e := range v.EdgesOut {
			var ok bool
			e.Target, ok = self.Vertexes[e.TargetID]
			if !ok {
				return fmt.Errorf("could not resolve Clout, bad TargetID: \"%s\"", e.TargetID)
			}

			e.Source = v
			e.Metadata = ard.EnsureMap(e.Metadata)
			e.Properties = ard.EnsureMap(e.Properties)

			e.Target.EdgesIn = append(e.Target.EdgesIn, e)
		}
	}
	return nil
}

func (self *Clout) Normalize() (*Clout, error) {
	s, err := format.EncodeYaml(self)
	if err != nil {
		return nil, err
	}
	return DecodeYaml(strings.NewReader(s))
}

//
// Entity
//

type Entity interface {
	GetMetadata() ard.Map
	GetProperties() ard.Map
}

//
// Vertex
//

type Vertex struct {
	Metadata   ard.Map `json:"metadata" yaml:"metadata"`
	ID         string  `json:"-" yaml:"-"`
	Properties ard.Map `json:"properties" yaml:"properties"`
	EdgesOut   Edges   `json:"edgesOut" yaml:"edgesOut"`
	EdgesIn    Edges   `json:"-" yaml:"-"`
}

func (self *Clout) NewVertex(id string) *Vertex {
	vertex := &Vertex{
		ID:         id,
		Metadata:   make(ard.Map),
		Properties: make(ard.Map),
		EdgesOut:   make(Edges, 0),
		EdgesIn:    make(Edges, 0),
	}
	self.Vertexes[id] = vertex
	return vertex
}

// Entity interface
func (self *Vertex) GetMetadata() ard.Map {
	return self.Metadata
}

// Entity interface
func (self *Vertex) GetProperties() ard.Map {
	return self.Properties
}

type Vertexes map[string]*Vertex

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
