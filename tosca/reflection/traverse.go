package reflection

import (
	"reflect"
)

type Traverser func(interface{}) bool

// Ignore fields tagged with "traverse:ignore" or "lookup"
func Traverse(entityPtr interface{}, traverse Traverser) {
	if !traverse(entityPtr) {
		return
	}

	if !IsPtrToStruct(reflect.TypeOf(entityPtr)) {
		return
	}

	value := reflect.ValueOf(entityPtr).Elem()

	for _, structField := range GetStructFields(value.Type()) {
		// Has traverse:"ignore" tag?
		traverseTag, ok := structField.Tag.Lookup("traverse")
		if ok && (traverseTag == "ignore") {
			continue
		}

		// Has "lookup" tag?
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
			for i := 0; i < length; i++ {
				element := field.Index(i)
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
