package tosca

import (
	"fmt"
	"reflect"
	"strings"
)

// From "name" tag
func GetEntityTypeName(type_ reflect.Type) string {
	numField := type_.NumField()
	for i := 0; i < numField; i++ {
		structField := type_.Field(i)
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
func GetKey(entityPtr interface{}) string {
	mappable, ok := entityPtr.(Mappable)
	if !ok {
		panic(fmt.Sprintf("entity does not implement \"Mappable\" interface: %T", entityPtr))
	}
	return mappable.GetKey()
}
