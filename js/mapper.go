package js

import (
	"reflect"
)

//
// goja.FieldNameMapper interface
//

var mapper Mapper

type Mapper struct{}

func (self Mapper) FieldName(t reflect.Type, f reflect.StructField) string {
	return ToJavaScriptStyle(f.Name)
}

func (self Mapper) MethodName(t reflect.Type, m reflect.Method) string {
	return ToJavaScriptStyle(m.Name)
}
