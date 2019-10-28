package js

import (
	"reflect"
	"unicode"
)

var mapper Mapper

type Mapper struct{}

// goja.FieldNameMapper interface
func (self Mapper) FieldName(t reflect.Type, f reflect.StructField) string {
	return ToJavaScriptStyle(f.Name)
}

// goja.FieldNameMapper interface
func (self Mapper) MethodName(t reflect.Type, m reflect.Method) string {
	return ToJavaScriptStyle(m.Name)
}

func ToJavaScriptStyle(name string) string {
	runes := []rune(name)
	length := len(runes)
	if (length > 0) && unicode.IsUpper(runes[0]) {
		if (length > 1) && unicode.IsUpper(runes[1]) {
			// If the second rune is also uppercase we'll keep the name as is
			return name
		}
		r := make([]rune, 1, length-1)
		r[0] = unicode.ToLower(runes[0])
		return string(append(r, runes[1:]...))
	}
	return name
}
