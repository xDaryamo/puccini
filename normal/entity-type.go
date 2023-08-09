package normal

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// EntityType
//

type EntityType struct {
	Name string `json:"-" yaml:"-"`

	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Parent      string            `json:"parent,omitempty" yaml:"parent,omitempty"`
}

func NewEntityType(name string) *EntityType {
	return &EntityType{
		Name:     name,
		Metadata: make(map[string]string),
	}
}

//
// EntityTypes
//

type EntityTypes map[string]*EntityType

func NewEntityTypes(names ...string) EntityTypes {
	entityTypes := make(EntityTypes)
	for _, name := range names {
		entityTypes[name] = NewEntityType(name)
	}
	return entityTypes
}

func GetHierarchyEntityTypes(hierarchy *parsing.Hierarchy) EntityTypes {
	entityTypes := make(EntityTypes)

	hierarchy.Range(func(entityPtr parsing.EntityPtr, parentEntityPtr parsing.EntityPtr) bool {
		entityType := NewEntityType(parsing.GetCanonicalName(entityPtr))

		if parentEntityPtr != nil {
			entityType.Parent = parsing.GetCanonicalName(parentEntityPtr)
		}

		entityType.Description, _ = parsing.GetDescription(entityPtr)

		if metadata, ok := parsing.GetMetadata(entityPtr); ok {
			for name, value := range metadata {
				// No need to include "canonical_name" metadata
				if name != "canonical_name" {
					entityType.Metadata[name] = value
				}
			}
		}

		entityTypes[entityType.Name] = entityType

		return true
	})

	return entityTypes
}

func GetEntityTypes(hierarchy *parsing.Hierarchy, entityPtr parsing.EntityPtr) (EntityTypes, bool) {
	if childHierarchy, ok := hierarchy.Find(entityPtr); ok {
		return GetHierarchyEntityTypes(childHierarchy), true
	}
	return nil, false
}
