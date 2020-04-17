package ard

import (
	"fmt"
	"time"
)

type Value = interface{}

// Note: This is just a convenient alias, *not* a type. An extra type would ensure more strictness but
// would make life more complicated than it needs to be. That said, if we *do* want to make this into a
// type, we need to make sure not to add any methods to the type, otherwise the goja JavaScript engine
// will treat it as a host object instead of a regular JavaScript dict object.
type List = []Value

// Note: This is just a convenient alias, *not* a type. An extra type would ensure more strictness but
// would make life more complicated than it needs to be. That said, if we *do* want to make this into a
// type, we need to make sure not to add any methods to the type, otherwise the goja JavaScript engine
// will treat it as a host object instead of a regular JavaScript dict object.
type Map = map[Value]Value

type StringMap = map[string]Value

// We will use YAML schema names:
// https://yaml.org/spec/1.2/spec.html#Schema

type TypeValidator = func(Value) bool

var TypeValidators = map[string]TypeValidator{
	// Failsafe schema: https://yaml.org/spec/1.2/spec.html#id2802346
	"!!map": IsMap,
	"!!seq": IsList,
	"!!str": IsString,

	// JSON schema: https://yaml.org/spec/1.2/spec.html#id2803231
	"!!bool":  IsBool,
	"!!int":   IsInteger,
	"!!float": IsFloat,

	// Other schemas: https://yaml.org/spec/1.2/spec.html#id2805770
	"!!null":      IsNull,
	"!!timestamp": IsTime,
}

func TypeName(value Value) string {
	switch value.(type) {
	case Map:
		return "!!map"
	case List:
		return "!!seq"
	case string:
		return "!!str"
	case bool:
		return "!!bool"
	case int64, int32, int16, int8, int:
		return "!!int"
	case float64, float32:
		return "!!float"
	case nil:
		return "!!null"
	case time.Time:
		return "!!timestamp"
	default:
		return fmt.Sprintf("%T", value)
	}
}

var TypeZeroes = map[string]Value{
	"!!map":       make(Map),
	"!!seq":       List{},
	"!!str":       "",
	"!!bool":      false,
	"!!int":       int(0),       // YAML parser returns int
	"!!float":     float64(0.0), // YAML parser returns float64
	"!!null":      nil,
	"!!timestamp": time.Time{}, // YAML parser returns time.Time
}

// Map = map[interface{}]interface{}
func IsMap(value Value) bool {
	_, ok := value.(Map)
	return ok
}

// List = []interface{}
func IsList(value Value) bool {
	_, ok := value.(List)
	return ok
}

// string
func IsString(value Value) bool {
	_, ok := value.(string)
	return ok
}

// bool
func IsBool(value Value) bool {
	_, ok := value.(bool)
	return ok
}

// int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint
func IsInteger(value Value) bool {
	switch value.(type) {
	case int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint:
		return true
	}
	return false
}

// float64, float32
func IsFloat(value Value) bool {
	switch value.(type) {
	case float64, float32:
		return true
	}
	return false
}

func IsNull(value Value) bool {
	return value == nil
}

// time.Time
func IsTime(value Value) bool {
	_, ok := value.(time.Time)
	return ok
}
