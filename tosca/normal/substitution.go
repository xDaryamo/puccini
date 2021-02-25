package normal

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
)

//
// Substitution
//

type Substitution struct {
	ServiceTemplate *ServiceTemplate

	Type                string
	TypeMetadata        map[string]string
	CapabilityMappings  map[*NodeTemplate]*Capability
	RequirementMappings map[*NodeTemplate]string
	PropertyMappings    map[*NodeTemplate]string
	AttributeMappings   map[*NodeTemplate]string
	InterfaceMappings   map[*NodeTemplate]string
}

func (self *ServiceTemplate) NewSubstitution() *Substitution {
	substitutionMappings := &Substitution{
		ServiceTemplate:     self,
		TypeMetadata:        make(map[string]string),
		CapabilityMappings:  make(map[*NodeTemplate]*Capability),
		RequirementMappings: make(map[*NodeTemplate]string),
		PropertyMappings:    make(map[*NodeTemplate]string),
		AttributeMappings:   make(map[*NodeTemplate]string),
		InterfaceMappings:   make(map[*NodeTemplate]string),
	}
	self.Substitution = substitutionMappings
	return substitutionMappings
}

type MarshalableSubstitution struct {
	Type                string            `json:"type" yaml:"type"`
	TypeMetadata        map[string]string `json:"typeMetadata" yaml:"typeMetadata"`
	CapabilityMappings  map[string]string `json:"capabilityMappings" yaml:"capabilityMappings"`
	RequirementMappings map[string]string `json:"requirementMappings" yaml:"requirementMappings"`
	PropertyMappings    map[string]string `json:"propertyMappings" yaml:"propertyMappings"`
	AttributeMappings   map[string]string `json:"attributeMappings" yaml:"attributeMappings"`
	InterfaceMappings   map[string]string `json:"interfaceMappings" yaml:"interfaceMappings"`
}

func (self *Substitution) Marshalable() interface{} {
	capabilityMappings := make(map[string]string)
	for n, c := range self.CapabilityMappings {
		capabilityMappings[n.Name] = c.Name
	}

	requirementMappings := make(map[string]string)
	for n, r := range self.RequirementMappings {
		requirementMappings[n.Name] = r
	}

	propertyMappings := make(map[string]string)
	for n, p := range self.PropertyMappings {
		propertyMappings[n.Name] = p
	}

	attributeMappings := make(map[string]string)
	for n, a := range self.AttributeMappings {
		attributeMappings[n.Name] = a
	}

	interfaceMappings := make(map[string]string)
	for n, i := range self.InterfaceMappings {
		interfaceMappings[n.Name] = i
	}

	return &MarshalableSubstitution{
		Type:                self.Type,
		TypeMetadata:        self.TypeMetadata,
		CapabilityMappings:  capabilityMappings,
		RequirementMappings: requirementMappings,
		PropertyMappings:    propertyMappings,
		AttributeMappings:   attributeMappings,
		InterfaceMappings:   interfaceMappings,
	}
}

// json.Marshaler interface
func (self *Substitution) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Marshalable())
}

// yaml.Marshaler interface
func (self *Substitution) MarshalYAML() (interface{}, error) {
	return self.Marshalable(), nil
}

// cbor.Marshaler interface
func (self *Substitution) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(self.Marshalable())
}
