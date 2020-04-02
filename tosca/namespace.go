package tosca

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/tliron/puccini/common/terminal"
	"github.com/tliron/puccini/tosca/reflection"
)

type NameTransformer = func(string, interface{}) []string

//
// Namespace
//

type Namespace struct {
	namespace map[reflect.Type]map[string]interface{}
	lock      sync.RWMutex
}

func NewNamespace() *Namespace {
	return &Namespace{
		namespace: make(map[reflect.Type]map[string]interface{}),
	}
}

// From "namespace" tags
func NewNamespaceFor(entityPtr interface{}) *Namespace {
	self := NewNamespace()

	reflection.Traverse(entityPtr, func(entityPtr interface{}) bool {
		for _, field := range reflection.GetTaggedFields(entityPtr, "namespace") {
			if field.Kind() != reflect.String {
				panic(fmt.Sprintf("\"namespace\" tag can only be used on \"string\" field in struct: %T", entityPtr))
			}
			self.set(field.String(), entityPtr)
		}
		return true
	})

	return self
}

func (self *Namespace) Empty() bool {
	self.lock.RLock()
	defer self.lock.RUnlock()

	return len(self.namespace) == 0
}

func (self *Namespace) Range(f func(interface{}, interface{}) bool) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	for _, forType := range self.namespace {
		for _, entityPtr := range forType {
			if !f(forType, entityPtr) {
				return
			}
		}
	}
}

func (self *Namespace) Lookup(name string) (interface{}, bool) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	for _, forType := range self.namespace {
		if entityPtr, ok := forType[name]; ok {
			return entityPtr, true
		}
	}

	return nil, false
}

func (self *Namespace) LookupForType(name string, type_ reflect.Type) (interface{}, bool) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	if forType, ok := self.namespace[type_]; ok {
		entityPtr, ok := forType[name]
		return entityPtr, ok
	} else {
		return nil, false
	}
}

// If the name has already been set returns existing entityPtr, true
func (self *Namespace) Set(name string, entityPtr interface{}) (interface{}, bool) {
	self.lock.Lock()
	defer self.lock.Unlock()

	return self.set(name, entityPtr)
}

func (self *Namespace) set(name string, entityPtr interface{}) (interface{}, bool) {
	type_ := reflect.TypeOf(entityPtr)
	forType, ok := self.namespace[type_]
	if !ok {
		forType = make(map[string]interface{})
		self.namespace[type_] = forType
	}

	if existing, ok := forType[name]; ok {
		if existing != entityPtr {
			// We are trying to give the name to a different entity
			return existing, true
		}

		// We already have this entity at this name, and that's fine
		return nil, false
	}

	forType[name] = entityPtr

	return nil, false
}

func (self *Namespace) Merge(namespace *Namespace, nameTransformer NameTransformer) {
	if self == namespace {
		return
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	namespace.lock.RLock()
	defer namespace.lock.RUnlock()

	for type_, forType := range namespace.namespace {
		for name, entityPtr := range forType {
			var names []string

			if nameTransformer != nil {
				names = append(names, nameTransformer(name, entityPtr)...)
			} else {
				names = []string{name}
			}

			for _, name = range names {
				if existing, exists := self.set(name, entityPtr); exists {
					GetContext(entityPtr).ReportNameAmbiguous(type_.Elem(), name, entityPtr, existing)
				}
			}
		}
	}
}

// Print

func (self *Namespace) Print(indent int) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	// Sort type names
	var types TypesByName
	for type_ := range self.namespace {
		types = append(types, type_)
	}
	sort.Sort(types)

	nameIndent := indent + 1
	for _, type_ := range types {
		forType := self.namespace[type_]
		terminal.PrintIndent(indent)
		fmt.Fprintf(terminal.Stdout, "%s\n", terminal.ColorTypeName(type_.Elem().String()))

		// Sort names
		names := make([]string, len(forType))
		i := 0
		for name := range forType {
			names[i] = name
			i += 1
		}
		sort.Strings(names)

		for _, name := range names {
			terminal.PrintIndent(nameIndent)
			fmt.Fprintf(terminal.Stdout, "%s\n", name)
		}
	}
}

// sort.Interface

type TypesByName []reflect.Type

func (self TypesByName) Len() int {
	return len(self)
}

func (self TypesByName) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self TypesByName) Less(i, j int) bool {
	iName := self[i].Elem().String()
	jName := self[j].Elem().String()
	return strings.Compare(iName, jName) < 0
}
