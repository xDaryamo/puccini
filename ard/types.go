package ard

import (
	"fmt"
	"time"
)

// We will use YAML schema names:
// https://yaml.org/spec/1.2/spec.html#Schema

type TypeValidator = func(interface{}) bool

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
	"!!timestamp": IsTime,
}

func TypeName(value interface{}) string {
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
	case time.Time:
		return "!!timestamp"
	default:
		return fmt.Sprintf("%T", value)
	}
}

var TypeZeroes = map[string]interface{}{
	"!!map":       make(Map),
	"!!seq":       List{},
	"!!str":       "",
	"!!bool":      false,
	"!!int":       int(0),       // YAML parser returns int
	"!!float":     float64(0.0), // YAML parser returns float64
	"!!timestamp": time.Time{},  // YAML parser returns time.Time
}

// Map = map[interface{}]interface{}
func IsMap(value interface{}) bool {
	_, ok := value.(Map)
	return ok
}

// List = []interface{}
func IsList(value interface{}) bool {
	_, ok := value.(List)
	return ok
}

// string
func IsString(value interface{}) bool {
	_, ok := value.(string)
	return ok
}

// bool
func IsBool(value interface{}) bool {
	_, ok := value.(bool)
	return ok
}

// int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint
func IsInteger(value interface{}) bool {
	switch value.(type) {
	case int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint:
		return true
	}
	return false
}

// float64, float32
func IsFloat(value interface{}) bool {
	switch value.(type) {
	case float64, float32:
		return true
	}
	return false
}

// time.Time
func IsTime(value interface{}) bool {
	_, ok := value.(time.Time)
	return ok
}
