package clout

import (
	"github.com/tliron/puccini/ard"
)

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
