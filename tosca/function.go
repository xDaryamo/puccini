package tosca

// Note: This is conceptually part of the "tosca.normal" package, but must be separated due
// do circular import problems.

//
// Function
//

type Function struct {
	URL       string        `json:"url" yaml:"url"`
	Path      string        `json:"path" yaml:"path"`
	Location  string        `json:"location" yaml:"location"`
	Name      string        `json:"name" yaml:"name"`
	Arguments []interface{} `json:"arguments" yaml:"arguments"`
}

func NewFunction(url string, path string, location string, name string, arguments []interface{}) *Function {
	self := Function{
		URL:       url,
		Path:      path,
		Location:  location,
		Name:      name,
		Arguments: arguments,
	}
	return &self
}

func (self *Context) NewFunction(name string, arguments []interface{}) *Function {
	return NewFunction(self.URL.String(), self.Path.String(), self.Location(), name, arguments)
}
