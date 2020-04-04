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

	hierarchy.Range(func(entityPtr tosca.EntityPtr, parentEntityPtr tosca.EntityPtr) bool {
		type_ := NewType(tosca.GetCanonicalName(entityPtr))

		if parentEntityPtr != nil {
			type_.Parent = tosca.GetCanonicalName(parentEntityPtr)
		}

		type_.Description, _ = tosca.GetDescription(entityPtr)

		if metadata, ok := tosca.GetMetadata(entityPtr); ok {
			for name, value := range metadata {
				// No need to include "canonical_name" metadata
				if name != "canonical_name" {
					type_.Metadata[name] = value
				}
			}
		}

		types[type_.Name] = type_

		return true
	})

	return types
}

func GetTypes(hierarchy *tosca.Hierarchy, entityPtr tosca.EntityPtr) (Types, bool) {
	if childHierarchy, ok := hierarchy.Find(entityPtr); ok {
		return GetHierarchyTypes(childHierarchy), true
	}
	return nil, false
}
