package reflection

import (
	"fmt"
	"reflect"
	"sync"
)

// Includes fields "inherited" from anonymous struct pointer fields
// The order of field definition is important! Later fields will override previous fields
func GetStructFields(type_ reflect.Type) []reflect.StructField {
	if v, ok := structFieldsCache.Load(type_); ok {
		return v.([]reflect.StructField)
	}

	var structFields []reflect.StructField
	numField := type_.NumField()
	for i := 0; i < numField; i++ {
		structField := type_.Field(i)
		if structField.Anonymous && (structField.Type.Kind() == reflect.Ptr) {
			for _, structField = range GetStructFields(structField.Type.Elem()) {
				structFields = appendStructField(structFields, structField)
			}
		} else {
			structFields = appendStructField(structFields, structField)
		}
	}

	structFieldsCache.Store(type_, structFields)

	return structFields
}

func appendStructField(structFields []reflect.StructField, structField reflect.StructField) []reflect.StructField {
	found := false
	for index, f := range structFields {
		if f.Name == structField.Name {
			// Override
			structFields[index] = structField
			found = true
			break
		}
	}
	if !found {
		structFields = append(structFields, structField)
	}
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
