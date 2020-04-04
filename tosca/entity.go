package tosca

import (
	"fmt"
	"reflect"
	"strings"
)

type EntityPtr = interface{}

// From "name" tag
func GetEntityTypeName(type_ reflect.Type) string {
	fields := type_.NumField()
	for index := 0; index < fields; index++ {
		structField := type_.Field(index)
		if value, ok := structField.Tag.Lookup("name"); ok {
			return value
		}
	}
	return fmt.Sprintf("%s", type_)
}

//
// EntityPtrs
//

type EntityPtrs []interface{}

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
// Mappable
//

type Mappable interface {
	GetKey() string
}

// From Mappable interface
func GetKey(entityPtr EntityPtr) string {
	mappable, ok := entityPtr.(Mappable)
	if !ok {
		panic(fmt.Sprintf("entity does not implement \"Mappable\" interface: %T", entityPtr))
	}
	return mappable.GetKey()
}
