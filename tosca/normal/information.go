package normal

//
// TypeInformation
//

type TypeInformation struct {
	Name              string            `json:"name,omitempty" yaml:"name,omitempty"`
	Description       string            `json:"description,omitempty" yaml:"description,omitempty"`
	SchemaDescription string            `json:"schemaDescription,omitempty" yaml:"schemaDescription,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

func NewTypeInformation() *TypeInformation {
	return &TypeInformation{
		Metadata: make(map[string]string),
	}
}

func CopyTypeInformation(information *TypeInformation) *TypeInformation {
	if information != nil {
		self := TypeInformation{
			Name:              information.Name,
			Description:       information.Description,
			SchemaDescription: information.SchemaDescription,
		}
		if (information.Metadata != nil) && (len(information.Metadata) > 0) {
			self.Metadata = make(map[string]string)
			for key, value := range information.Metadata {
				self.Metadata[key] = value
			}
		}
		if !self.Empty() {
			return &self
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func (self *TypeInformation) Empty() bool {
	return (self.Name == "") &&
		(self.Description == "") &&
		(self.SchemaDescription == "") &&
		((self.Metadata == nil) || (len(self.Metadata) == 0))
}

func (self *TypeInformation) Prune() {
	if (self.Metadata != nil) && (len(self.Metadata) == 0) {
		self.Metadata = nil
	}
}

//
// Information
//

type Information struct {
	Description string           `json:"description,omitempty" yaml:"description,omitempty"`
	Definition  *TypeInformation `json:"definition,omitempty" yaml:"definition,omitempty"`
	Type        *TypeInformation `json:"type,omitempty" yaml:"type,omitempty"`
	Entry       *TypeInformation `json:"entry,omitempty" yaml:"entry,omitempty"`
	Key         *TypeInformation `json:"key,omitempty" yaml:"key,omitempty"`
	Value       *TypeInformation `json:"value,omitempty" yaml:"value,omitempty"`
	Properties  InformationMap   `json:"properties,omitempty" yaml:"properties,omitempty"`
}

func NewInformation() *Information {
	return &Information{
		Definition: NewTypeInformation(),
		Type:       NewTypeInformation(),
		Entry:      NewTypeInformation(),
		Key:        NewTypeInformation(),
		Value:      NewTypeInformation(),
		Properties: make(InformationMap),
	}
}

func CopyInformation(information *Information) *Information {
	if information != nil {
		self := Information{
			Description: information.Description,
			Definition:  CopyTypeInformation(information.Definition),
			Type:        CopyTypeInformation(information.Type),
			Entry:       CopyTypeInformation(information.Entry),
			Key:         CopyTypeInformation(information.Key),
			Value:       CopyTypeInformation(information.Value),
			Properties:  CopyInformationMap(information.Properties),
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

func (self *Information) Empty() bool {
	return (self.Description == "") &&
		((self.Definition == nil) || self.Definition.Empty()) &&
		((self.Type == nil) || self.Type.Empty()) &&
		((self.Entry == nil) || self.Entry.Empty()) &&
		((self.Key == nil) || self.Key.Empty()) &&
		((self.Value == nil) || self.Value.Empty()) &&
		((self.Properties == nil) || self.Properties.Empty())
}

func (self *Information) Prune() {
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
	if self.Properties != nil {
		if self.Properties.Prune(); self.Properties.Empty() {
			self.Properties = nil
		}
	}
}

//
// InformationMap
//

type InformationMap map[string]*Information

func CopyInformationMap(informationMap InformationMap) InformationMap {
	if (informationMap == nil) || (len(informationMap) == 0) {
		return nil
	}
	self := make(InformationMap)
	for key, information := range informationMap {
		self[key] = CopyInformation(information)
	}
	if !self.Empty() {
		return self
	} else {
		return nil
	}
}

func (self InformationMap) Empty() bool {
	for _, information := range self {
		if !information.Empty() {
			return false
		}
	}
	return true
}

func (self InformationMap) Prune() {
	for key, information := range self {
		if information.Prune(); information.Empty() {
			delete(self, key)
		}
	}
}
