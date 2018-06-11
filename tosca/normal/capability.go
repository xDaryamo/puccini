package normal

//
// Capability
//

type Capability struct {
	NodeTemplate *NodeTemplate  `json:"-" yaml:"-"`
	Name         string         `json:"-" yaml:"-"`
	Description  string         `json:"description" yaml:"description"`
	Types        Types          `json:"types" yaml:"types"`
	Properties   Constrainables `json:"properties" yaml:"properties"`
	Attributes   Constrainables `json:"attributes" yaml:"attributes"`
}

func (self *NodeTemplate) NewCapability(name string) *Capability {
	capability := &Capability{
		NodeTemplate: self,
		Name:         name,
		Types:        make(Types),
		Properties:   make(Constrainables),
		Attributes:   make(Constrainables),
	}
	self.Capabilities[name] = capability
	return capability
}

//
// Capabilities
//

type Capabilities map[string]*Capability
