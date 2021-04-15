package normal

//
// Substitution
//

type Substitution struct {
	ServiceTemplate *ServiceTemplate

	Type                string
	TypeMetadata        map[string]string
	CapabilityMappings  Mappings
	RequirementMappings Mappings
	PropertyMappings    Mappings
	AttributeMappings   Mappings
	InterfaceMappings   Mappings
}

func (self *ServiceTemplate) NewSubstitution() *Substitution {
	substitutionMappings := &Substitution{
		ServiceTemplate:     self,
		TypeMetadata:        make(map[string]string),
		CapabilityMappings:  make(Mappings),
		RequirementMappings: make(Mappings),
		PropertyMappings:    make(Mappings),
		AttributeMappings:   make(Mappings),
		InterfaceMappings:   make(Mappings),
	}
	self.Substitution = substitutionMappings
	return substitutionMappings
}
