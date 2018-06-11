package normal

//
// Policy
//

type Policy struct {
	ServiceTemplate     *ServiceTemplate `json:"-" yaml:"-"`
	Name                string           `json:"-" yaml:"-"`
	Description         string           `json:"description" yaml:"description"`
	Types               Types            `json:"types" yaml:"types"`
	Properties          Constrainables   `json:"properties" yaml:"properties"`
	GroupTargets        []*Group         `json:"-" yaml:"-"`
	NodeTemplateTargets []*NodeTemplate  `json:"-" yaml:"-"`
}

func (self *ServiceTemplate) NewPolicy(name string) *Policy {
	policy := &Policy{
		ServiceTemplate:     self,
		Name:                name,
		Types:               make(Types),
		Properties:          make(Constrainables),
		GroupTargets:        make([]*Group, 0),
		NodeTemplateTargets: make([]*NodeTemplate, 0),
	}
	self.Policies[name] = policy
	return policy
}

//
// Policies
//

type Policies map[string]*Policy
