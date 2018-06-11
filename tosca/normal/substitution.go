package normal

import (
	"encoding/json"
)

//
// Substitution
//

type Substitution struct {
	ServiceTemplate     *ServiceTemplate
	Type                string
	TypeMetadata        map[string]string
	CapabilityMappings  map[*NodeTemplate]*Capability
	RequirementMappings map[*NodeTemplate]string
}

func (self *ServiceTemplate) NewSubstitution() *Substitution {
	substitutionMappings := &Substitution{
		ServiceTemplate:     self,
		TypeMetadata:        make(map[string]string),
		CapabilityMappings:  make(map[*NodeTemplate]*Capability),
		RequirementMappings: make(map[*NodeTemplate]string),
	}
	self.Substitution = substitutionMappings
	return substitutionMappings
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

	return &struct {
		Type                string            `json:"type" yaml:"type"`
		TypeMetadata        map[string]string `json:"typeMetadata" yaml:"typeMetadata"`
		CapabilityMappings  map[string]string `json:"capabilityMappings" yaml:"capabilityMappings"`
		RequirementMappings map[string]string `json:"requirementMappings" yaml:"requirementMappings"`
	}{
		Type:                self.Type,
		TypeMetadata:        self.TypeMetadata,
		CapabilityMappings:  capabilityMappings,
		RequirementMappings: requirementMappings,
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
