package normal

import (
	"github.com/tliron/puccini/tosca"
)

//
// Type
//

type Type struct {
	Name string `json:"-" yaml:"-"`

	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Parent      string            `json:"parent,omitempty" yaml:"parent,omitempty"`
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

		type_.Description, _ = tosca.GetDescription(h.EntityPtr)

		if metadata, ok := tosca.GetMetadata(h.EntityPtr); ok {
			for name, value := range metadata {
				// No need to include "canonical_name" metadata
				if name != "canonical_name" {
					type_.Metadata[name] = value
				}
			}
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
