package normal

//
// Operation
//

type Operation struct {
	Interface     *Interface     `json:"-" yaml:"-"`
	PolicyTrigger *PolicyTrigger `json:"-" yaml:"-"`
	Name          string         `json:"-" yaml:"-"`

	Description    string   `json:"description" yaml:"description"`
	Implementation string   `json:"implementation" yaml:"implementation"`
	Dependencies   []string `json:"dependencies" yaml:"dependencies"`
	Inputs         Values   `json:"inputs" yaml:"inputs"`
	Outputs        Values   `json:"outputs" yaml:"outputs"`
	Timeout        int64    `json:"timeout" yaml:"timeout"`
	Host           string   `json:"host,omitempty" yaml:"host,omitempty"`
}

func (self *Interface) NewOperation(name string) *Operation {
	operation := &Operation{
		Interface:    self,
		Name:         name,
		Dependencies: make([]string, 0),
		Inputs:       make(Values),
		Outputs:      make(Values),
		Timeout:      -1,
	}
	self.Operations[name] = operation
	return operation
}

func (self *PolicyTrigger) NewOperation() *Operation {
	self.Operation = &Operation{
		PolicyTrigger: self,
		Dependencies:  make([]string, 0),
		Inputs:        make(Values),
		Outputs:       make(Values),
	}
	return self.Operation
}

//
// Operations
//

type Operations map[string]*Operation
