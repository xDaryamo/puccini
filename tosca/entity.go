package tosca

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/reflection"
	"github.com/tliron/kutil/url"
)

//
// EntityPtr
//

type EntityPtr = any

//
// EntityPtrs
//

type EntityPtrs []EntityPtr

// sort.Interface

func (self EntityPtrs) Len() int {
	return len(self)
}

func (self EntityPtrs) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self EntityPtrs) Less(i, j int) bool {
	iName := GetContext(self[i]).Path.String()
	jName := GetContext(self[j]).Path.String()
	return strings.Compare(iName, jName) < 0
}

//
// EntitySet
//

type EntityPtrSet map[EntityPtr]struct{}

func (self EntityPtrSet) Add(entityPtr EntityPtr) {
	self[entityPtr] = struct{}{}
}

func (self EntityPtrSet) Contains(entityPtr EntityPtr) bool {
	_, ok := self[entityPtr]
	return ok
}

//
// PreReadable
//

type PreReadable interface {
	PreRead()
}

// From PreReadable interface
func PreRead(entityPtr EntityPtr) bool {
	if preReadable, ok := entityPtr.(PreReadable); ok {
		preReadable.PreRead()
		return true
	} else {
		return false
	}
}

//
// Importer
//

type Importer interface {
	GetImportSpecs() []*ImportSpec
}

// From Importer interface
func GetImportSpecs(entityPtr EntityPtr) []*ImportSpec {
	if importer, ok := entityPtr.(Importer); ok {
		return importer.GetImportSpecs()
	} else {
		return nil
	}
}

//
// ImportSpec
//

type ImportSpec struct {
	URL             url.URL
	NameTransformer NameTransformer
	Implicit        bool
}

//
// Hierarchical
//

type Hierarchical interface {
	GetParent() EntityPtr
}

// From Hierarchical interface
func GetParent(entityPtr EntityPtr) (EntityPtr, bool) {
	if hierarchical, ok := entityPtr.(Hierarchical); ok {
		parentPtr := hierarchical.GetParent()
		if reflect.ValueOf(parentPtr).IsNil() {
			parentPtr = nil
		}
		return parentPtr, true
	} else {
		return nil, false
	}
}

//
// Inherits
//

type Inherits interface {
	Inherit()
}

// From Inherits interface
func Inherit(entityPtr EntityPtr) bool {
	if inherits, ok := entityPtr.(Inherits); ok {
		inherits.Inherit()
		return true
	} else {
		return false
	}
}

//
// Renderable
//

type Renderable interface {
	Render()
}

// From Renderable interface
func Render(entityPtr EntityPtr) bool {
	if renderable, ok := entityPtr.(Renderable); ok {
		renderable.Render()
		return true
	} else {
		return false
	}
}

//
// Mappable
//

type Mappable interface {
	GetKey() string
}

// From Mappable interface
func GetKey(entityPtr EntityPtr) string {
	if mappable, ok := entityPtr.(Mappable); ok {
		return mappable.GetKey()
	} else {
		panic(fmt.Sprintf("entity does not implement \"Mappable\" interface: %T", entityPtr))
	}
}

//
// HasInputs
//

type HasInputs interface {
	SetInputs(map[string]ard.Value)
}

// From HasInputs interface
func SetInputs(entityPtr EntityPtr, inputs map[string]ard.Value) bool {
	var done bool

	if inputs == nil {
		return false
	}

	reflection.TraverseEntities(entityPtr, false, func(entityPtr EntityPtr) bool {
		if hasInputs, ok := entityPtr.(HasInputs); ok {
			hasInputs.SetInputs(inputs)
			done = true

			// Only one entity should implement the interface
			return false
		}
		return true
	})

	return done
}

//
// HasMetadata
//

const (
	METADATA_TYPE                    = "puccini.type"
	METADATA_CONVERTER               = "puccini.converter"
	METADATA_COMPARER                = "puccini.comparer"
	METADATA_QUIRKS                  = "puccini.quirks"
	METADATA_DATA_TYPE_PREFIX        = "puccini.data-type:"
	METADATA_SCRIPTLET_PREFIX        = "puccini.scriptlet:"
	METADATA_SCRIPTLET_IMPORT_PREFIX = "puccini.scriptlet.import:"

	METADATA_CANONICAL_NAME    = "tosca.canonical-name"
	METADATA_NORMATIVE         = "tosca.normative"
	METADATA_FUNCTION_PREFIX   = "tosca.function."
	METADATA_CONSTRAINT_PREFIX = "tosca.constraint."
)

type HasMetadata interface {
	GetDescription() (string, bool)
	GetMetadata() (map[string]string, bool) // should return a copy
	SetMetadata(name string, value string) bool
}

// From HasMetadata interface
func GetDescription(entityPtr EntityPtr) (string, bool) {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		return hasMetadata.GetDescription()
	} else {
		return "", false
	}
}

// From HasMetadata interface
func GetMetadata(entityPtr EntityPtr) (map[string]string, bool) {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		return hasMetadata.GetMetadata()
	} else {
		return nil, false
	}
}

// From HasMetadata interface
func SetMetadata(entityPtr EntityPtr, name string, value string) bool {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		hasMetadata.SetMetadata(name, value)
		return true
	} else {
		return false
	}
}

func GetDataTypeMetadata(metadata map[string]string) map[string]string {
	dataTypeMetadata := make(map[string]string)
	if metadata != nil {
		for key, value := range metadata {
			if strings.HasPrefix(key, METADATA_DATA_TYPE_PREFIX) {
				dataTypeMetadata[key[len(METADATA_DATA_TYPE_PREFIX):]] = value
			}
		}
	}
	return dataTypeMetadata
}
