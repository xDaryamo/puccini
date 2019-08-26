package tosca

// Note: This is conceptually part of the "tosca.normal" package, but must be separated to
// avoid circular imports.

//
// FunctionCall
//

type FunctionCall struct {
	URL       string        `json:"url" yaml:"url"`
	Path      string        `json:"path" yaml:"path"`
	Location  string        `json:"location" yaml:"location"`
	Name      string        `json:"name" yaml:"name"`
	Arguments []interface{} `json:"arguments" yaml:"arguments"`
}

func NewFunctionCall(url string, path string, location string, name string, arguments []interface{}) *FunctionCall {
	return &FunctionCall{
		URL:       url,
		Path:      path,
		Location:  location,
		Name:      name,
		Arguments: arguments,
	}
}

func (self *Context) NewFunctionCall(name string, arguments []interface{}) *FunctionCall {
	return NewFunctionCall(self.URL.String(), self.Path.String(), self.Location(), name, arguments)
}
