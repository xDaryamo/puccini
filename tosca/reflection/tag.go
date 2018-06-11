package reflection

import (
	"reflect"
)

//
// FieldTag
//

type FieldTag struct {
	FieldName string
	Tag       string
}

func GetTaggedFields(entityPtr interface{}, name string) []reflect.Value {
	var fields []reflect.Value
	entity := reflect.ValueOf(entityPtr).Elem()
	for _, tag := range GetFieldTagsForValue(entity, name) {
		field := entity.FieldByName(tag.FieldName)
		fields = append(fields, field)
	}
	return fields
}

func GetFieldTagsForValue(value reflect.Value, name string) []FieldTag {
	return GetFieldTagsForType(value.Type(), name)
}

func GetFieldTagsForType(type_ reflect.Type, name string) []FieldTag {
	var tags []FieldTag
	for _, structField := range GetStructFields(type_) {
		if value, ok := structField.Tag.Lookup(name); ok {
			tags = append(tags, FieldTag{structField.Name, value})
		}
	}
	return tags
}
