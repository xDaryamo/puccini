package tosca

import (
	"fmt"
	"strings"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca/problems"
	"github.com/tliron/puccini/url"
	"github.com/tliron/yamlkeys"
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
	Parent             *Context
	Name               string
	Path               ard.Path
	URL                url.URL
	Data               interface{} // ARD
	Locator            ard.Locator
	Namespace          Namespace
	ScriptletNamespace ScriptletNamespace
	Hierarchy          *Hierarchy
	Problems           *problems.Problems
	Quirks             []string
	Grammar            Grammar
	ReadOverrides      map[string]string
}

func NewContext(problems *problems.Problems, quirks []string) *Context {
	return &Context{
		Namespace:          make(Namespace),
		ScriptletNamespace: make(ScriptletNamespace),
		Hierarchy:          &Hierarchy{},
		Problems:           problems,
		Quirks:             quirks,
	}
}

func (self *Context) NewImportContext(url_ url.URL) *Context {
	return &Context{
		Name:               self.Name,
		Path:               self.Path,
		URL:                url_,
		Namespace:          make(Namespace),
		ScriptletNamespace: make(ScriptletNamespace),
		Hierarchy:          &Hierarchy{},
		Problems:           self.Problems,
		Quirks:             self.Quirks,
		Grammar:            self.Grammar,
	}
}

func (self *Context) Clone(data interface{}) *Context {
	return &Context{
		Parent:             self.Parent,
		Name:               self.Name,
		Path:               self.Path,
		URL:                self.URL,
		Data:               data,
		Locator:            self.Locator,
		Namespace:          self.Namespace,
		ScriptletNamespace: self.ScriptletNamespace,
		Hierarchy:          self.Hierarchy,
		Problems:           self.Problems,
		Quirks:             self.Quirks,
		Grammar:            self.Grammar,
	}
}

func (self *Context) GetAncestor(generation int) *Context {
	if generation == 0 {
		return self
	} else if generation == 1 {
		return self.Parent
	} else if self.Parent != nil {
		return self.Parent.GetAncestor(generation - 1)
	} else {
		return nil
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

func (self *Context) HasQuirk(quirk string) bool {
	for _, q := range self.Quirks {
		if q == quirk {
			return true
		}
	}
	return false
}

func (self *Context) Location() string {
	if self.Locator != nil {
		if r, c, ok := self.Locator.Locate(self.Path...); ok {
			return fmt.Sprintf("%d,%d", r, c)
		}
	}
	return ""
}

//
// Child contexts
//

func (self *Context) FieldChild(name interface{}, data interface{}) *Context {
	nameString := yamlkeys.KeyString(name) // complex keys would be stringified
	return &Context{
		Parent:             self,
		Name:               nameString,
		Path:               append(self.Path, ard.NewFieldPathElement(nameString)),
		URL:                self.URL,
		Data:               data,
		Locator:            self.Locator,
		Namespace:          self.Namespace,
		ScriptletNamespace: self.ScriptletNamespace,
		Hierarchy:          self.Hierarchy,
		Problems:           self.Problems,
		Quirks:             self.Quirks,
		Grammar:            self.Grammar,
	}
}

func (self *Context) GetFieldChild(name string) (*Context, bool) {
	if self.ValidateType("map") {
		if data, ok := self.Data.(ard.Map)[name]; ok {
			return self.FieldChild(name, data), true
		}
	}
	return nil, false
}

func (self *Context) GetRequiredFieldChild(name string) (*Context, bool) {
	if context, ok := self.GetFieldChild(name); ok {
		return context, true
	} else {
		self.FieldChild(name, nil).ReportFieldMissing()
		return nil, false
	}
}

func (self *Context) FieldChildren() []*Context {
	var children []*Context
	for name, data := range self.Data.(ard.Map) {
		children = append(children, self.FieldChild(name, data))
	}
	return children
}

func (self *Context) MapChild(name interface{}, data interface{}) *Context {
	nameString := strings.ReplaceAll(yamlkeys.KeyString(name), "\n", "¶") // complex keys would be stringified

	return &Context{
		Parent:             self,
		Name:               nameString,
		Path:               append(self.Path, ard.NewMapPathElement(nameString)),
		URL:                self.URL,
		Data:               data,
		Locator:            self.Locator,
		Namespace:          self.Namespace,
		ScriptletNamespace: self.ScriptletNamespace,
		Hierarchy:          self.Hierarchy,
		Problems:           self.Problems,
		Quirks:             self.Quirks,
		Grammar:            self.Grammar,
	}
}

func (self *Context) ListChild(index int, data interface{}) *Context {
	return &Context{
		Parent:             self,
		Name:               fmt.Sprintf("%d", index),
		Path:               append(self.Path, ard.NewListPathElement(index)),
		URL:                self.URL,
		Data:               data,
		Locator:            self.Locator,
		Namespace:          self.Namespace,
		ScriptletNamespace: self.ScriptletNamespace,
		Hierarchy:          self.Hierarchy,
		Problems:           self.Problems,
		Quirks:             self.Quirks,
		Grammar:            self.Grammar,
	}
}

func (self *Context) SequencedListChild(index int, name string, data interface{}) *Context {
	return &Context{
		Parent:             self,
		Name:               name,
		Path:               append(self.Path, ard.NewListPathElement(index)),
		URL:                self.URL,
		Data:               data,
		Locator:            self.Locator,
		Namespace:          self.Namespace,
		ScriptletNamespace: self.ScriptletNamespace,
		Hierarchy:          self.Hierarchy,
		Problems:           self.Problems,
		Quirks:             self.Quirks,
		Grammar:            self.Grammar,
	}
}
