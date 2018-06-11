package tosca

import (
	"fmt"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca/problems"
	"github.com/tliron/puccini/url"
)

//
// Contextual
//

type Contextual interface {
	GetContext() *Context
}

// From Contextual interface
func GetContext(entityPtr interface{}) *Context {
	contextual, ok := entityPtr.(Contextual)
	if !ok {
		panic(fmt.Sprintf("entity does not implement \"Contextual\" interface: %T", entityPtr))
	}
	return contextual.GetContext()
}

//
// Context
//

type Context struct {
	Parent          *Context
	Name            string
	Path            string
	URL             url.URL
	Data            interface{}
	Namespace       Namespace
	ScriptNamespace ScriptNamespace
	Hierarchy       *Hierarchy
	Problems        *problems.Problems
}

func NewContext(problems *problems.Problems) Context {
	return Context{
		Namespace:       make(Namespace),
		ScriptNamespace: make(ScriptNamespace),
		Hierarchy:       &Hierarchy{},
		Problems:        problems,
	}
}

func (self *Context) Is(typeNames ...string) bool {
	valid := false
	for _, typeName := range typeNames {
		typeValidator, ok := PrimitiveTypeValidators[typeName]
		if !ok {
			panic(fmt.Sprintf("unsupported field type: %s", typeName))
		}
		if typeValidator(self.Data) {
			valid = true
			break
		}
	}
	return valid
}

//
// Child contexts
//

func (self *Context) FieldChild(name string, data interface{}) *Context {
	var path string
	if self.Path == "" {
		path = name
	} else {
		path = fmt.Sprintf("%s.%s", self.Path, name)
	}

	return &Context{
		Parent:          self,
		Name:            name,
		Path:            path,
		URL:             self.URL,
		Data:            data,
		Namespace:       self.Namespace,
		ScriptNamespace: self.ScriptNamespace,
		Hierarchy:       self.Hierarchy,
		Problems:        self.Problems,
	}
}

func (self *Context) RequiredFieldChild(name string) (*Context, bool) {
	if !self.ValidateType("map") {
		return nil, false
	}

	data, ok := self.Data.(ard.Map)[name]
	if !ok {
		self.FieldChild(name, nil).ReportFieldMissing()
		return nil, false
	}

	return self.FieldChild(name, data), true
}

func (self *Context) FieldChildren() []*Context {
	var children []*Context
	for name, data := range self.Data.(ard.Map) {
		children = append(children, self.FieldChild(name, data))
	}
	return children
}

func (self *Context) MapChild(name string, data interface{}) *Context {
	return &Context{
		Parent:          self,
		Name:            name,
		Path:            fmt.Sprintf("%s['%s']", self.Path, name),
		URL:             self.URL,
		Data:            data,
		Namespace:       self.Namespace,
		ScriptNamespace: self.ScriptNamespace,
		Hierarchy:       self.Hierarchy,
		Problems:        self.Problems,
	}
}

func (self *Context) ListChild(index int, data interface{}) *Context {
	return &Context{
		Parent:          self,
		Name:            fmt.Sprintf("%d", index),
		Path:            fmt.Sprintf("%s[%d]", self.Path, index),
		URL:             self.URL,
		Data:            data,
		Namespace:       self.Namespace,
		ScriptNamespace: self.ScriptNamespace,
		Hierarchy:       self.Hierarchy,
		Problems:        self.Problems,
	}
}

func (self *Context) SequencedListChild(index int, name string, data interface{}) *Context {
	return &Context{
		Parent:          self,
		Name:            name,
		Path:            fmt.Sprintf("%s[%d]", self.Path, index),
		URL:             self.URL,
		Data:            data,
		Namespace:       self.Namespace,
		ScriptNamespace: self.ScriptNamespace,
		Hierarchy:       self.Hierarchy,
		Problems:        self.Problems,
	}
}

func (self *Context) Import(url_ url.URL) *Context {
	return &Context{
		Name:            self.Name,
		Path:            self.Path,
		URL:             url_,
		Namespace:       make(Namespace),
		ScriptNamespace: make(ScriptNamespace),
		Hierarchy:       &Hierarchy{},
		Problems:        self.Problems,
	}
}

func (self *Context) WithData(data interface{}) *Context {
	return &Context{
		Parent:          self.Parent,
		Name:            self.Name,
		Path:            self.Path,
		URL:             self.URL,
		Data:            data,
		Namespace:       self.Namespace,
		ScriptNamespace: self.ScriptNamespace,
		Hierarchy:       self.Hierarchy,
		Problems:        self.Problems,
	}
}
