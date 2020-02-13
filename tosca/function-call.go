package tosca

// Note: This is conceptually part of the "tosca.normal" package, but must be separated to
// avoid circular imports.

//
// FunctionCall
//

type FunctionCall struct {
	Name      string        `json:"name" yaml:"name"`
	Arguments []interface{} `json:"arguments" yaml:"arguments"`
	URL       string        `json:"url,omitempty" yaml:"url,omitempty"`
	Row       int           `json:"row,omitempty" yaml:"row,omitempty"`
	Column    int           `json:"column,omitempty" yaml:"column,omitempty"`
	Path      string        `json:"path,omitempty" yaml:"path,omitempty"`
}

func NewFunctionCall(name string, arguments []interface{}, url string, row int, column int, path string) *FunctionCall {
	return &FunctionCall{
		Name:      name,
		Arguments: arguments,
		URL:       url,
		Row:       row,
		Column:    column,
		Path:      path,
	}
}

func (self *Context) NewFunctionCall(name string, arguments []interface{}) *FunctionCall {
	row, column := self.GetLocation()
	return NewFunctionCall(name, arguments, self.URL.String(), row, column, self.Path.String())
}
