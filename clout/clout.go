package clout

import (
	"encoding/json"
	"errors"
	"fmt"

	"strings"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/format"
)

const Version = "1.0"

//
// Clout
//

type Clout struct {
	Version    string        `json:"version" yaml:"version"`
	Metadata   ard.StringMap `json:"metadata" yaml:"metadata"`
	Properties ard.StringMap `json:"properties" yaml:"properties" `
	Vertexes   Vertexes      `json:"vertexes" yaml:"vertexes"`
}

func NewClout() *Clout {
	return &Clout{
		Version:    Version,
		Metadata:   make(ard.StringMap),
		Properties: make(ard.StringMap),
		Vertexes:   make(Vertexes),
	}
}

type MarshalableCloutStringMaps Clout

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
	// JavaScript requires keys to be strings, so we would lose complex keys
	return json.Marshal(self.MarshalableStringMaps())
}

func (self *Clout) Resolve() error {
	if self.Version == "" {
		return errors.New("no Clout \"Version\"")
	}
	if self.Version != Version {
		return fmt.Errorf("unsupported Clout version: %q", self.Version)
	}

	for key, v := range self.Vertexes {
		v.Clout = self
		v.ID = key

		for _, e := range v.EdgesOut {
			var ok bool
			if e.Target, ok = self.Vertexes[e.TargetID]; !ok {
				return fmt.Errorf("could not resolve Clout, bad TargetID: %q", e.TargetID)
			}

			e.Source = v
			e.Target.EdgesIn = append(e.Target.EdgesIn, e)
		}
	}
	return nil
}

func (self *Clout) Normalize() (*Clout, error) {
	// TODO: not very efficient
	if code, err := format.EncodeYAML(self, " ", false); err == nil {
		return Read(strings.NewReader(code), "yaml")
	} else {
		return nil, err
	}
}
