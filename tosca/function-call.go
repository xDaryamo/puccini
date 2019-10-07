package tosca

// Note: This is conceptually part of the "tosca.normal" package, but must be separated to
// avoid circular imports.

//
// FunctionCall
//

type FunctionCall struct {
	Name      string        `json:"name" yaml:"name"`
	Arguments []interface{} `json:"arguments" yaml:"arguments"`
	URL       string        `json:"url" yaml:"url"`
	Location  string        `json:"location" yaml:"location"`
	Path      string        `json:"path" yaml:"path"`
}

func NewFunctionCall(name string, arguments []interface{}, url string, location string, path string) *FunctionCall {
	return &FunctionCall{
		Name:      name,
		Arguments: arguments,
		URL:       url,
		Location:  location,
		Path:      path,
	}
}

func (self *Context) NewFunctionCall(name string, arguments []interface{}) *FunctionCall {
	return NewFunctionCall(name, arguments, self.URL.String(), self.Location(), self.Path.String())
}
