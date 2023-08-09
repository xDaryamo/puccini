package normal

//
// Substitution
//

type Substitution struct {
	ServiceTemplate *ServiceTemplate `json:"-" yaml:"-"`

	Type               string            `json:"type" yaml:"type"`
	TypeMetadata       map[string]string `json:"typeMetadata" yaml:"typeMetadata"`
	InputPointers      Pointers          `json:"inputPointers" yaml:"inputPointers"`
	CapabilityPointers Pointers          `json:"capabilityPointers" yaml:"capabilityPointers"`
	RequirementPointer Pointers          `json:"requirementPointers" yaml:"requirementPointers"`
	PropertyPointers   Pointers          `json:"propertyPointers" yaml:"propertyPointers"`
	PropertyValues     Values            `json:"propertyValues" yaml:"propertyValues"`
	AttributePointers  Pointers          `json:"attributePointers" yaml:"attributePointers"`
	InterfacePointers  Pointers          `json:"interfacePointers" yaml:"interfacePointers"`
}

func (self *ServiceTemplate) NewSubstitution() *Substitution {
	substitution := &Substitution{
		ServiceTemplate:    self,
		TypeMetadata:       make(map[string]string),
		InputPointers:      make(Pointers),
		CapabilityPointers: make(Pointers),
		RequirementPointer: make(Pointers),
		PropertyPointers:   make(Pointers),
		PropertyValues:     make(Values),
		AttributePointers:  make(Pointers),
		InterfacePointers:  make(Pointers),
	}
	self.Substitution = substitution
	return substitution
}
