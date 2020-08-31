package tosca

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/problems"
	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/yamlkeys"
)

//
// Contextual
//

type Contextual interface {
	GetContext() *Context
}

// From Contextual interface
func GetContext(entityPtr EntityPtr) *Context {
	if contextual, ok := entityPtr.(Contextual); ok {
		return contextual.GetContext()
	} else {
		panic(fmt.Sprintf("entity does not implement \"Contextual\" interface: %T", entityPtr))
	}
}

func GetCanonicalName(entityPtr EntityPtr) string {
	if metadata, ok := GetMetadata(entityPtr); ok {
		if canonicalName, ok := metadata["canonical_name"]; ok {
			return canonicalName
		}
	}

	context := GetContext(entityPtr)
	canonicalNamespace := context.GetCanonicalNamespace()
	if canonicalNamespace != nil {
		return fmt.Sprintf("%s::%s", *canonicalNamespace, context.Name)
	} else {
		return context.Name
	}
}

//
// Context
//

type Context struct {
	Parent             *Context
	Name               string
	Path               ard.Path
	URL                urlpkg.URL
	Data               ard.Value
	Locator            ard.Locator
	CanonicalNamespace *string
	Namespace          *Namespace
	ScriptletNamespace *ScriptletNamespace
	Hierarchy          *Hierarchy
	Problems           *problems.Problems
	Quirks             Quirks
	Grammar            *Grammar
	ReadTagOverrides   map[string]string
}

func NewContext(quirks Quirks) *Context {
	return &Context{
		Namespace:          NewNamespace(),
		ScriptletNamespace: NewScriptletNamespace(),
		Hierarchy:          NewHierarchy(),
		Problems:           new(problems.Problems),
		Quirks:             quirks,
	}
}

func (self *Context) NewImportContext(url urlpkg.URL) *Context {
	return &Context{
		Name:               self.Name,
		Path:               self.Path,
		URL:                url,
		CanonicalNamespace: self.CanonicalNamespace,
		Namespace:          NewNamespace(),
		ScriptletNamespace: NewScriptletNamespace(),
		Hierarchy:          NewHierarchy(),
		Problems:           new(problems.Problems),
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

func (self *Context) GetCanonicalNamespace() *string {
	if self.CanonicalNamespace != nil {
		return self.CanonicalNamespace
	} else if self.Parent != nil {
		return self.Parent.GetCanonicalNamespace()
	}
	return nil
}

func (self *Context) HasQuirk(quirk Quirk) bool {
	return self.Quirks.Has(quirk)
}

func (self *Context) SetReadTag(fieldName string, tag string) {
	if self.ReadTagOverrides == nil {
		self.ReadTagOverrides = make(map[string]string)
	}
	self.ReadTagOverrides[fieldName] = tag
}

func (self *Context) GetLocation() (int, int) {
	if self.Locator != nil {
		if row, column, ok := self.Locator.Locate(self.Path...); ok {
			return row, column
		}
	}
	return -1, -1
}

func (self *Context) Is(typeNames ...ard.TypeName) bool {
	for _, typeName := range typeNames {
		if typeValidator, ok := ard.TypeValidators[typeName]; ok {
			if typeValidator(self.Data) {
				return true
			}
		} else {
			panic(fmt.Sprintf("unsupported field type: %s", typeName))
		}
	}
	return false
}

//
// Child contexts
//

func (self *Context) Clone(data ard.Value) *Context {
	return &Context{
		Parent:             self.Parent,
		Name:               self.Name,
		Path:               self.Path,
		URL:                self.URL,
		Data:               data,
		Locator:            self.Locator,
		CanonicalNamespace: self.CanonicalNamespace,
		Namespace:          self.Namespace,
		ScriptletNamespace: self.ScriptletNamespace,
		Hierarchy:          self.Hierarchy,
		Problems:           self.Problems,
		Quirks:             self.Quirks,
		Grammar:            self.Grammar,
	}
}

func (self *Context) FieldChild(name ard.Value, data ard.Value) *Context {
	nameString := yamlkeys.KeyString(name) // complex keys would be stringified
	return &Context{
		Parent:             self,
		Name:               nameString,
		Path:               self.Path.AppendField(nameString),
		URL:                self.URL,
		Data:               data,
		Locator:            self.Locator,
		CanonicalNamespace: self.CanonicalNamespace,
		Namespace:          self.Namespace,
		ScriptletNamespace: self.ScriptletNamespace,
		Hierarchy:          self.Hierarchy,
		Problems:           self.Problems,
		Quirks:             self.Quirks,
		Grammar:            self.Grammar,
	}
}

func (self *Context) GetFieldChild(name string) (*Context, bool) {
	if self.ValidateType(ard.TypeMap) {
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

func (self *Context) MapChild(name ard.Value, data ard.Value) *Context {
	nameString := strings.ReplaceAll(yamlkeys.KeyString(name), "\n", "¶") // complex keys would be stringified
	return &Context{
		Parent:             self,
		Name:               nameString,
		Path:               self.Path.AppendMap(nameString),
		URL:                self.URL,
		Data:               data,
		Locator:            self.Locator,
		CanonicalNamespace: self.CanonicalNamespace,
		Namespace:          self.Namespace,
		ScriptletNamespace: self.ScriptletNamespace,
		Hierarchy:          self.Hierarchy,
		Problems:           self.Problems,
		Quirks:             self.Quirks,
		Grammar:            self.Grammar,
	}
}

func (self *Context) ListChild(index int, data ard.Value) *Context {
	return &Context{
		Parent:             self,
		Name:               strconv.FormatInt(int64(index), 10),
		Path:               self.Path.AppendList(index),
		URL:                self.URL,
		Data:               data,
		Locator:            self.Locator,
		CanonicalNamespace: self.CanonicalNamespace,
		Namespace:          self.Namespace,
		ScriptletNamespace: self.ScriptletNamespace,
		Hierarchy:          self.Hierarchy,
		Problems:           self.Problems,
		Quirks:             self.Quirks,
		Grammar:            self.Grammar,
	}
}

func (self *Context) SequencedListChild(index int, name string, data ard.Value) *Context {
	return &Context{
		Parent:             self,
		Name:               name,
		Path:               self.Path.AppendSequencedList(index),
		URL:                self.URL,
		Data:               data,
		Locator:            self.Locator,
		CanonicalNamespace: self.CanonicalNamespace,
		Namespace:          self.Namespace,
		ScriptletNamespace: self.ScriptletNamespace,
		Hierarchy:          self.Hierarchy,
		Problems:           self.Problems,
		Quirks:             self.Quirks,
		Grammar:            self.Grammar,
	}
}
