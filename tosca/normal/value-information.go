package normal

//
// ValueInformation
//

type ValueInformation struct {
	Description string              `json:"description,omitempty" yaml:"description,omitempty"`
	Definition  *TypeInformation    `json:"definition,omitempty" yaml:"definition,omitempty"`
	Type        *TypeInformation    `json:"type,omitempty" yaml:"type,omitempty"`
	Entry       *TypeInformation    `json:"entry,omitempty" yaml:"entry,omitempty"`
	Key         *TypeInformation    `json:"key,omitempty" yaml:"key,omitempty"`
	Value       *TypeInformation    `json:"value,omitempty" yaml:"value,omitempty"`
	Fields      ValueInformationMap `json:"fields,omitempty" yaml:"fields,omitempty"`
}

func NewValueInformation() *ValueInformation {
	return &ValueInformation{
		Definition: NewTypeInformation(),
		Type:       NewTypeInformation(),
		Entry:      NewTypeInformation(),
		Key:        NewTypeInformation(),
		Value:      NewTypeInformation(),
		Fields:     make(ValueInformationMap),
	}
}

func CopyValueInformation(information *ValueInformation) *ValueInformation {
	if information != nil {
		self := ValueInformation{
			Description: information.Description,
			Definition:  CopyTypeInformation(information.Definition),
			Type:        CopyTypeInformation(information.Type),
			Entry:       CopyTypeInformation(information.Entry),
			Key:         CopyTypeInformation(information.Key),
			Value:       CopyTypeInformation(information.Value),
			Fields:      CopyValueInformationMap(information.Fields),
		}
		if self.Prune(); !self.Empty() {
			return &self
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func (self *ValueInformation) Empty() bool {
	return (self.Description == "") &&
		((self.Definition == nil) || self.Definition.Empty()) &&
		((self.Type == nil) || self.Type.Empty()) &&
		((self.Entry == nil) || self.Entry.Empty()) &&
		((self.Key == nil) || self.Key.Empty()) &&
		((self.Value == nil) || self.Value.Empty()) &&
		((self.Fields == nil) || self.Fields.Empty())
}

func (self *ValueInformation) Prune() {
	if self.Definition != nil {
		if self.Definition.Prune(); self.Definition.Empty() {
			self.Definition = nil
		}
	}
	if self.Type != nil {
		if self.Type.Prune(); self.Type.Empty() {
			self.Type = nil
		}
	}
	if self.Entry != nil {
		if self.Entry.Prune(); self.Entry.Empty() {
			self.Entry = nil
		}
	}
	if self.Key != nil {
		if self.Key.Prune(); self.Key.Empty() {
			self.Key = nil
		}
	}
	if self.Value != nil {
		if self.Value.Prune(); self.Value.Empty() {
			self.Value = nil
		}
	}
	if self.Fields != nil {
		if self.Fields.Prune(); self.Fields.Empty() {
			self.Fields = nil
		}
	}
}

//
// ValueInformationMap
//

type ValueInformationMap map[string]*ValueInformation

func CopyValueInformationMap(informationMap ValueInformationMap) ValueInformationMap {
	if (informationMap == nil) || (len(informationMap) == 0) {
		return nil
	}
	self := make(ValueInformationMap)
	for key, information := range informationMap {
		self[key] = CopyValueInformation(information)
	}
	if !self.Empty() {
		return self
	} else {
		return nil
	}
}

func (self ValueInformationMap) Empty() bool {
	for _, information := range self {
		if !information.Empty() {
			return false
		}
	}
	return true
}

func (self ValueInformationMap) Prune() {
	for key, information := range self {
		if information.Prune(); information.Empty() {
			delete(self, key)
		}
	}
}
