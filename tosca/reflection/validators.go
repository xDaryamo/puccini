package reflection

import (
	"time"

	"github.com/tliron/puccini/ard"
)

//
// Validators
//

type TypeValidator func(interface{}) bool

// bool
func IsBool(value interface{}) bool {
	_, ok := value.(bool)
	return ok
}

// int, int64, int32
func IsInt(value interface{}) bool {
	switch value.(type) {
	case int, int64, int32:
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

// ard.List = []interface{}
func IsSliceOfStruct(value interface{}) bool {
	_, ok := value.(ard.List)
	return ok
}

// ard.Map = map[interface{}]interface{}
func IsMap(value interface{}) bool {
	_, ok := value.(ard.Map)
	return ok
}

// *string
func IsPtrToString(value interface{}) bool {
	_, ok := value.(*string)
	return ok
}

// *int64
func IsPtrToInt64(value interface{}) bool {
	_, ok := value.(*int64)
	return ok
}

// *float64
func IsPtrToFloat64(value interface{}) bool {
	_, ok := value.(*float64)
	return ok
}

// *bool
func IsPtrToBool(value interface{}) bool {
	_, ok := value.(*bool)
	return ok
}

// *[]string
func IsPtrToSliceOfString(value interface{}) bool {
	_, ok := value.(*[]string)
	return ok
}

// *map[string]string
func IsPtrToMapOfStringToString(value interface{}) bool {
	_, ok := value.(*map[string]string)
	return ok
}
