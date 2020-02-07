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
	Parent   string            `json:"parent,omitempty" yaml:"parent,omitempty"`
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
	h := hierarchy
	for (h != nil) && (h.EntityPtr != nil) {
		type_ := NewType(tosca.GetCanonicalName(h.EntityPtr))

		if (h.Parent != nil) && (h.Parent.EntityPtr != nil) {
			type_.Parent = tosca.GetCanonicalName(h.Parent.EntityPtr)
		}

		if metadata, ok := tosca.GetMetadata(h.EntityPtr); ok {
			type_.Metadata = metadata
		}

		types[type_.Name] = type_

		h = h.Parent
	}
	return types
}

func GetTypes(hierarchy *tosca.Hierarchy, entityPtr interface{}) (Types, bool) {
	if childHierarchy, ok := hierarchy.Find(entityPtr); ok {
		return GetHierarchyTypes(childHierarchy), true
	}
	return nil, false
}
