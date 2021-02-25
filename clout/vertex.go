package clout

import (
	"encoding/json"

	"github.com/tliron/kutil/ard"
)

//
// Vertex
//

type Vertex struct {
	Clout      *Clout        `yaml:"-" cbor:"-"`
	ID         string        `yaml:"-" cbor:"-"`
	Metadata   ard.StringMap `yaml:"metadata" cbor:"metadata"`
	Properties ard.StringMap `yaml:"properties" cbor:"properties"`
	EdgesOut   Edges         `yaml:"edgesOut" cbor:"edgesOut"`
	EdgesIn    Edges         `yaml:"-" cbor:"-"`
}

func (self *Clout) NewVertex(id string) *Vertex {
	vertex := &Vertex{
		Clout:      self,
		ID:         id,
		Metadata:   make(ard.StringMap),
		Properties: make(ard.StringMap),
		EdgesOut:   make(Edges, 0),
		EdgesIn:    make(Edges, 0),
	}
	self.Vertexes[id] = vertex
	return vertex
}

func (self *Vertex) Remove() {
	delete(self.Clout.Vertexes, self.ID)
}

// Entity interface
func (self *Vertex) GetMetadata() ard.StringMap {
	return self.Metadata
}

// Entity interface
func (self *Vertex) GetProperties() ard.StringMap {
	return self.Properties
}

type MarshalableVertexStringMaps struct {
	Metadata   ard.StringMap `json:"metadata"`
	Properties ard.StringMap `json:"properties"`
	EdgesOut   Edges         `json:"edgesOut"`
}

func (self *Vertex) MarshalableStringMaps() interface{} {
	return &MarshalableVertexStringMaps{
		Metadata:   ard.EnsureStringMaps(self.Metadata),
		Properties: ard.EnsureStringMaps(self.Properties),
		EdgesOut:   self.EdgesOut,
	}
}

// json.Marshaler interface
func (self *Vertex) MarshalJSON() ([]byte, error) {
	// JavaScript requires keys to be strings, so we would lose complex keys
	return json.Marshal(self.MarshalableStringMaps())
}

//
// Vertexes
//

type Vertexes map[string]*Vertex
