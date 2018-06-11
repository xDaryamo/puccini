package reflection

import (
	"fmt"
	"reflect"
	"sync"
)

// Includes fields "inherited" from anonymous struct pointer fields
func GetStructFields(type_ reflect.Type) []reflect.StructField {
	if v, ok := structFieldsCache.Load(type_); ok {
		return v.([]reflect.StructField)
	}

	var structFields []reflect.StructField
	numField := type_.NumField()
	for i := 0; i < numField; i++ {
		structField := type_.Field(i)
		if structField.Anonymous && (structField.Type.Kind() == reflect.Ptr) {
			structFields = append(structFields, GetStructFields(structField.Type.Elem())...)
		} else {
			structFields = append(structFields, structField)
		}
	}

	structFieldsCache.Store(type_, structFields)

	return structFields
}

var structFieldsCache sync.Map

func GetReferredField(entity reflect.Value, referenceFieldName string, referredFieldName string) (reflect.Value, reflect.Value, bool) {
	referenceField := entity.FieldByName(referenceFieldName)
	if !referenceField.IsValid() {
		panic(fmt.Sprintf("tag refers to unknown field \"%s\" in struct: %s", referenceFieldName, entity.Type()))
	}
	if referenceField.Type().Kind() != reflect.Ptr {
		panic(fmt.Sprintf("tag refers to non-pointer field \"%s\" in struct: %s", referenceFieldName, entity.Type()))
	}

	if referenceField.IsNil() {
		// Reference is empty
		return referenceField, reflect.Value{}, false
	}

	referredField := referenceField.Elem().FieldByName(referredFieldName)

	if !referredField.IsValid() {
		panic(fmt.Sprintf("tag's field name \"%s\" not found in the entity referred to by \"%s\" in struct: %s", referredFieldName, referenceFieldName, entity.Type()))
	}

	return referenceField, referredField, true
}
