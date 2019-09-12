package clout

import (
	"encoding/json"
	"errors"
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
	Version    string   `yaml:"version"`
	Metadata   ard.Map  ` yaml:"metadata"`
	Properties ard.Map  ` yaml:"properties"`
	Vertexes   Vertexes `yaml:"vertexes"`
}

func NewClout() *Clout {
	return &Clout{
		Version:    Version,
		Metadata:   make(ard.Map),
		Properties: make(ard.Map),
		Vertexes:   make(Vertexes),
	}
}

type MarshalableCloutStringMaps struct {
	Version    string        `json:"version"`
	Metadata   ard.StringMap `json:"metadata"`
	Properties ard.StringMap `json:"properties"`
	Vertexes   Vertexes      `json:"vertexes"`
}

func (self *Clout) MarshalableStringMaps() interface{} {
	return &MarshalableCloutStringMaps{
		Version:    self.Version,
		Metadata:   ard.EnsureStringMaps(self.Metadata),
		Properties: ard.EnsureStringMaps(self.Properties),
		Vertexes:   self.Vertexes,
	}
}

// json.Marshaler interface
func (self *Clout) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.MarshalableStringMaps())
}

func (self *Clout) Resolve() error {
	if self.Version == "" {
		return errors.New("no Clout \"Version\"")
	}
	if self.Version != Version {
		return fmt.Errorf("unsupported Clout version: \"%s\"", self.Version)
	}

	// TODO: do we need these?
	self.Metadata = ard.EnsureMaps(self.Metadata)
	self.Properties = ard.EnsureMaps(self.Properties)

	for key, v := range self.Vertexes {
		v.Clout = self
		v.ID = key
		// TODO: do we need these?
		v.Metadata = ard.EnsureMaps(v.Metadata)
		v.Properties = ard.EnsureMaps(v.Properties)

		for _, e := range v.EdgesOut {
			var ok bool
			if e.Target, ok = self.Vertexes[e.TargetID]; !ok {
				return fmt.Errorf("could not resolve Clout, bad TargetID: \"%s\"", e.TargetID)
			}

			e.Source = v
			// TODO: do we need these?
			e.Metadata = ard.EnsureMaps(e.Metadata)
			e.Properties = ard.EnsureMaps(e.Properties)

			e.Target.EdgesIn = append(e.Target.EdgesIn, e)
		}
	}
	return nil
}

func (self *Clout) Normalize() (*Clout, error) {
	return self, nil
	if s, err := format.EncodeYaml(self, " "); err == nil {
		return DecodeYaml(strings.NewReader(s))
	} else {
		return nil, err
	}
}
