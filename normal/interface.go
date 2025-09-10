package normal

//
// Interface
//

type Interface struct {
	NodeTemplate *NodeTemplate `json:"-" yaml:"-"`
	Group        *Group        `json:"-" yaml:"-"`
	Relationship *Relationship `json:"-" yaml:"-"`
	Name         string        `json:"-" yaml:"-"`

	Description   string            `json:"description" yaml:"description"`
	Metadata      map[string]string `json:"metadata" yaml:"metadata"`
	Types         EntityTypes       `json:"types" yaml:"types"`
	Inputs        Values            `json:"inputs" yaml:"inputs"`
	Operations    Operations        `json:"operations" yaml:"operations"`
	Notifications Notifications     `json:"notifications" yaml:"notifications"`
}

func (self *NodeTemplate) NewInterface(name string) *Interface {
	interface_ := &Interface{
		NodeTemplate:  self,
		Name:          name,
		Metadata:      make(map[string]string),
		Types:         make(EntityTypes),
		Inputs:        make(Values),
		Operations:    make(Operations),
		Notifications: make(Notifications),
	}
	self.Interfaces[name] = interface_
	return interface_
}

func (self *Group) NewInterface(name string) *Interface {
	interface_ := &Interface{
		Group:         self,
		Name:          name,
		Metadata:      make(map[string]string),
		Types:         make(EntityTypes),
		Inputs:        make(Values),
		Operations:    make(Operations),
		Notifications: make(Notifications),
	}
	self.Interfaces[name] = interface_
	return interface_
}

func (self *Relationship) NewInterface(name string) *Interface {
	interface_ := &Interface{
		Relationship:  self,
		Name:          name,
		Metadata:      make(map[string]string),
		Types:         make(EntityTypes),
		Inputs:        make(Values),
		Operations:    make(Operations),
		Notifications: make(Notifications),
	}
	self.Interfaces[name] = interface_
	return interface_
}

//
// Interfaces
//

type Interfaces map[string]*Interface
