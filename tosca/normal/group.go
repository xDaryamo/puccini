package normal

//
// Group
//

type Group struct {
	ServiceTemplate *ServiceTemplate `json:"-" yaml:"-"`
	Name            string           `json:"-" yaml:"-"`

	Metadata    map[string]string `json:"metadata" yaml:"metadata"`
	Description string            `json:"description" yaml:"description"`
	Types       Types             `json:"types" yaml:"types"`
	Properties  Constrainables    `json:"properties" yaml:"properties"`
	Interfaces  Interfaces        `json:"interfaces" yaml:"interfaces"`

	Members []*NodeTemplate `json:"-" yaml:"-"`
}

func (self *ServiceTemplate) NewGroup(name string) *Group {
	group := &Group{
		ServiceTemplate: self,
		Name:            name,
		Metadata:        make(map[string]string),
		Types:           make(Types),
		Properties:      make(Constrainables),
		Interfaces:      make(Interfaces),
		Members:         make([]*NodeTemplate, 0),
	}
	self.Groups[name] = group
	return group
}

//
// Groups
//

type Groups map[string]*Group
