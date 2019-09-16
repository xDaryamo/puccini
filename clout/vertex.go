package clout

import (
	"encoding/json"

	"github.com/tliron/puccini/ard"
)

//
// Vertex
//

type Vertex struct {
	Clout      *Clout        `yaml:"-"`
	ID         string        `yaml:"-"`
	Metadata   ard.StringMap `yaml:"metadata"`
	Properties ard.StringMap `yaml:"properties"`
	EdgesOut   Edges         `yaml:"edgesOut"`
	EdgesIn    Edges         `yaml:"-"`
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
