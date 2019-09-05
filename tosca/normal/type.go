package normal

import (
	"github.com/tliron/puccini/tosca"
)

//
// Type
//

type Type struct {
	Name string `json:"-" yaml:"-"`

	Metadata map[string]string `json:"metadata" yaml:"metadata"`
}

func NewType(name string) *Type {
	return &Type{
		Name:     name,
		Metadata: make(map[string]string),
	}
}

//
// Types
//

type Types map[string]*Type

func NewTypes(names ...string) Types {
	types := make(Types)
	for _, name := range names {
		types[name] = NewType(name)
	}
	return types
}

func GetHierarchyTypes(hierarchy *tosca.Hierarchy) Types {
	types := make(Types)
	n := hierarchy
	for (n != nil) && (n.EntityPtr != nil) {
		name := n.GetContext().Name
		type_ := NewType(name)
		if metadata, ok := GetMetadata(n.EntityPtr); ok {
			type_.Metadata = metadata
		}
		types[name] = type_
		n = n.Parent
	}
	return types
}

func GetTypes(hierarchy *tosca.Hierarchy, entityPtr interface{}) (Types, bool) {
	if childHierarchy, ok := hierarchy.Find(entityPtr); ok {
		return GetHierarchyTypes(childHierarchy), true
	}
	return nil, false
}
