package ard

import (
	"fmt"
	"time"
)

// We will use YAML type names:
// https://yaml.org/type/

type TypeValidator = func(interface{}) bool

var TypeValidators = map[string]TypeValidator{
	"!!bool":      IsBool,
	"!!int":       IsInt,
	"!!float":     IsFloat,
	"!!str":       IsString,
	"!!timestamp": IsTime,
	"!!seq":       IsList,
	"!!map":       IsMap,
}

var TypeZeroes = map[string]interface{}{
	"!!bool":      false,
	"!!int":       int(0),       // YAML parser returns int
	"!!float":     float64(0.0), // YAML parser returns float64
	"!!str":       "",
	"!!timestamp": time.Time{}, // YAML parser returns time.Time
	"!!seq":       List{},
	"!!map":       make(Map),
}

func TypeName(value interface{}) string {
	switch value.(type) {
	case bool:
		return "!!bool"
	case int64, int32, int16, int8, int:
		return "!!int"
	case float64, float32:
		return "!!float"
	case string:
		return "!!str"
	case time.Time:
		return "!!timestamp"
	case List:
		return "!!seq"
	case Map:
		return "!!map"
	default:
		return fmt.Sprintf("%T", value)
	}
}

// bool
func IsBool(value interface{}) bool {
	_, ok := value.(bool)
	return ok
}

// int64, int32, int16, int8, int
func IsInt(value interface{}) bool {
	switch value.(type) {
	case int64, int32, int16, int8, int:
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

// string
func IsString(value interface{}) bool {
	_, ok := value.(string)
	return ok
}

// time.Time
func IsTime(value interface{}) bool {
	_, ok := value.(time.Time)
	return ok
}

// List = []interface{}
func IsList(value interface{}) bool {
	_, ok := value.(List)
	return ok
}

// Map = map[interface{}]interface{}
func IsMap(value interface{}) bool {
	_, ok := value.(Map)
	return ok
}
