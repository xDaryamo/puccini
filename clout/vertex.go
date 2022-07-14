package clout

import (
	"encoding/json"

	"github.com/tliron/kutil/ard"
)

//
// Vertex
//

type Vertex struct {
	Clout      *Clout        `json:"-" yaml:"-"`
	ID         string        `json:"-" yaml:"-"`
	Metadata   ard.StringMap `json:"metadata" yaml:"metadata"`
	Properties ard.StringMap `json:"properties" yaml:"properties"`
	EdgesOut   Edges         `json:"edgesOut" yaml:"edgesOut"`
	EdgesIn    Edges         `json:"-" yaml:"-"`
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

func (self *Vertex) MarshalableStringMaps() any {
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

func (self *Vertex) copy(toArd bool) (*Vertex, error) {
	vertex := Vertex{
		ID:      self.ID,
		EdgesIn: make(Edges, 0),
	}

	var err error
	if vertex.EdgesOut, err = self.EdgesOut.copy(toArd); err == nil {
		if toArd {
			if metadata, err := ard.NormalizeStringMapsCopyToARD(self.Metadata); err == nil {
				if properties, err := ard.NormalizeStringMapsCopyToARD(self.Properties); err == nil {
					vertex.Metadata = metadata.(ard.StringMap)
					vertex.Properties = properties.(ard.StringMap)
				} else {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			vertex.Metadata = ard.SimpleCopy(self.Metadata).(ard.StringMap)
			vertex.Properties = ard.SimpleCopy(self.Properties).(ard.StringMap)
		}
	} else {
		return nil, err
	}

	return &vertex, nil
}

//
// Vertexes
//

type Vertexes map[string]*Vertex

func (self Vertexes) copy(toArd bool) (Vertexes, error) {
	vertexes := make(Vertexes)
	var err error
	for id, vertex := range self {
		if vertexes[id], err = vertex.copy(toArd); err != nil {
			return nil, err
		}
	}
	return vertexes, nil
}
