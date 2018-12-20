package clout

import (
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
		return errors.New("no Clout \"Version\"")
	}
	if self.Version != Version {
		return fmt.Errorf("unsupported Clout version: \"%s\"", self.Version)
	}

	self.Metadata = ard.EnsureMap(self.Metadata)
	self.Properties = ard.EnsureMap(self.Properties)

	for key, v := range self.Vertexes {
		v.Clout = self
		v.ID = key
		v.Metadata = ard.EnsureMap(v.Metadata)
		v.Properties = ard.EnsureMap(v.Properties)

		for _, e := range v.EdgesOut {
			var ok bool
			if e.Target, ok = self.Vertexes[e.TargetID]; !ok {
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
	if s, err := format.EncodeYaml(self); err == nil {
		return DecodeYaml(strings.NewReader(s))
	} else {
		return nil, err
	}
}
