package tosca

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/tosca/reflection"
)

type NameTransformer func(string, interface{}) []string

//
// Namespace
//

type Namespace map[reflect.Type]map[string]interface{}

// From "namespace" tags
func NewNamespace(entityPtr interface{}) Namespace {
	namespace := make(Namespace)

	reflection.Traverse(entityPtr, func(entityPtr interface{}) bool {
		for _, field := range reflection.GetTaggedFields(entityPtr, "namespace") {
			if field.Kind() != reflect.String {
				panic(fmt.Sprintf("\"namespace\" tag can only be used on \"string\" field in struct: %T", entityPtr))
			}
			namespace.Set(field.String(), entityPtr)
		}
		return true
	})

	return namespace
}

func (self Namespace) Lookup(name string) (interface{}, bool) {
	for _, forType := range self {
		if entityPtr, ok := forType[name]; ok {
			return entityPtr, true
		}
	}
	return nil, false
}

func (self Namespace) LookupForType(name string, type_ reflect.Type) (interface{}, bool) {
	if forType, ok := self[type_]; ok {
		entityPtr, ok := forType[name]
		return entityPtr, ok
	} else {
		return nil, false
	}
}

// If the name has already been set returns existing entityPtr, true
func (self Namespace) Set(name string, entityPtr interface{}) (interface{}, bool) {
	type_ := reflect.TypeOf(entityPtr)
	forType, ok := self[type_]
	if !ok {
		forType = make(map[string]interface{})
		self[type_] = forType
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

func (self Namespace) Merge(namespace Namespace, nameTransformer NameTransformer) {
	for type_, forType := range namespace {
		for name, entityPtr := range forType {
			var names []string

			if nameTransformer != nil {
				names = append(names, nameTransformer(name, entityPtr)...)
			} else {
				names = []string{name}
			}

			for _, name = range names {
				existing, exists := self.Set(name, entityPtr)
				if exists {
					GetContext(entityPtr).ReportNameAmbiguous(type_.Elem(), name, entityPtr, existing)
				}
			}
		}
	}
}

// Print

func (self Namespace) Print(indent int) {
	// Sort type names
	var types TypesByName
	for type_ := range self {
		types = append(types, type_)
	}
	sort.Sort(types)

	nameIndent := indent + 1
	for _, type_ := range types {
		forType := self[type_]
		format.PrintIndent(indent)
		fmt.Fprintf(format.Stdout, "%s\n", format.ColorTypeName(type_.Elem().String()))

		// Sort names
		names := make([]string, len(forType))
		i := 0
		for name := range forType {
			names[i] = name
			i += 1
		}
		sort.Strings(names)

		for _, name := range names {
			format.PrintIndent(nameIndent)
			fmt.Fprintf(format.Stdout, "%s\n", name)
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
