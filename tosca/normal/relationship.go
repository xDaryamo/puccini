package normal

//
// Relationship
//

type Relationship struct {
	Requirement *Requirement `json:"-" yaml:"-"`
	Name        string       `json:"-" yaml:"-"`

	Metadata    map[string]string `json:"metadata" yaml:"metadata"`
	Description string            `json:"description" yaml:"description"`
	Types       EntityTypes       `json:"types" yaml:"types"`
	Properties  Values            `json:"properties" yaml:"properties"`
	Attributes  Values            `json:"attributes" yaml:"attributes"`
	Interfaces  Interfaces        `json:"interfaces" yaml:"interfaces"`
}

func (self *Requirement) NewRelationship() *Relationship {
	relationship := &Relationship{
		Requirement: self,
		Metadata:    make(map[string]string),
		Types:       make(EntityTypes),
		Properties:  make(Values),
		Attributes:  make(Values),
		Interfaces:  make(Interfaces),
	}
	self.Relationship = relationship
	return relationship
}

//
// Relationships
//

type Relationships []*Relationship
