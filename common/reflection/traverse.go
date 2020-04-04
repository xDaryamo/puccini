package reflection

import (
	"reflect"
)

type Traverser func(interface{}) bool

// Ignore fields tagged with "traverse:ignore" or "lookup"
func Traverse(object interface{}, traverse Traverser) {
	if !traverse(object) {
		return
	}

	if !IsPtrToStruct(reflect.TypeOf(object)) {
		return
	}

	value := reflect.ValueOf(object).Elem()

	for _, structField := range GetStructFields(value.Type()) {
		// Has traverse:"ignore" tag?
		traverseTag, ok := structField.Tag.Lookup("traverse")
		if ok && (traverseTag == "ignore") {
			continue
		}

		// Ignore if has "lookup" tag
		if _, ok = structField.Tag.Lookup("lookup"); ok {
			continue
		}

		field := value.FieldByName(structField.Name)
		fieldType := field.Type()
		if IsPtrToStruct(fieldType) && !field.IsNil() {
			// Compatible with *interface{}
			Traverse(field.Interface(), traverse)
		} else if IsSliceOfPtrToStruct(fieldType) {
			// Compatible with []*interface{}
			length := field.Len()
			for index := 0; index < length; index++ {
				element := field.Index(index)
				Traverse(element.Interface(), traverse)
			}
		} else if IsMapOfStringToPtrToStruct(fieldType) {
			// Compatible with map[string]*interface{}
			for _, mapKey := range field.MapKeys() {
				element := field.MapIndex(mapKey)
				Traverse(element.Interface(), traverse)
			}
		}
	}
}
