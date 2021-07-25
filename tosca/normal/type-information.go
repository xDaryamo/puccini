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
