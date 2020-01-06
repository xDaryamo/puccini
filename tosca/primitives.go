package tosca

import (
	"fmt"
	"time"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca/reflection"
)

var PrimitiveTypeValidators = map[string]reflection.TypeValidator{
	"boolean":     reflection.IsBool,
	"integer":     reflection.IsInt,
	"float":       reflection.IsFloat,
	"string":      reflection.IsString,
	"!!timestamp": reflection.IsTime, // With the "!!" prefix so it won't conflict with TOSCA "timestamp" type name
	"list":        reflection.IsSliceOfStruct,
	"map":         reflection.IsMap,
}

var PrimitiveTypeZeroes = map[string]interface{}{
	"boolean":     false,
	"integer":     0,
	"float":       0.0,
	"string":      "",
	"!!timestamp": time.Time{},
	"list":        ard.List{},
	"map":         make(ard.Map),
}

func PrimitiveTypeName(value interface{}) string {
	switch value.(type) {
	case bool:
		return "boolean"
	case int: // YAML parser returns int
		return "integer"
	case float64: // YAML parser returns float64
		return "float"
	case time.Time: // YAML parser returns time.Time
		return "!!timestamp"
	case ard.List:
		return "list"
	case ard.Map:
		return "map"
	default:
		return fmt.Sprintf("%T", value)
	}
}
