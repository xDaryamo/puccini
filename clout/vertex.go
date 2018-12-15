package clout

import (
	"github.com/tliron/puccini/ard"
)

//
// Vertex
//

type Vertex struct {
	Clout      *Clout  `json:"-" yaml:"-"`
	ID         string  `json:"-" yaml:"-"`
	Metadata   ard.Map `json:"metadata" yaml:"metadata"`
	Properties ard.Map `json:"properties" yaml:"properties"`
	EdgesOut   Edges   `json:"edgesOut" yaml:"edgesOut"`
	EdgesIn    Edges   `json:"-" yaml:"-"`
}

func (self *Clout) NewVertex(id string) *Vertex {
	vertex := &Vertex{
		Clout:      self,
		ID:         id,
		Metadata:   make(ard.Map),
		Properties: make(ard.Map),
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
func (self *Vertex) GetMetadata() ard.Map {
	return self.Metadata
}

// Entity interface
func (self *Vertex) GetProperties() ard.Map {
	return self.Properties
}

//
// Vertexes
//

type Vertexes map[string]*Vertex
