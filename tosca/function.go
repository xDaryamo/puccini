package tosca

//
// Function
//

type Function struct {
	Path      string        `json:"path" yaml:"path"`
	Name      string        `json:"name" yaml:"name"`
	Arguments []interface{} `json:"arguments" yaml:"arguments"`
}

func NewFunction(path string, name string, arguments []interface{}) *Function {
	self := Function{
		Path:      path,
		Name:      name,
		Arguments: arguments,
	}
	return &self
}
