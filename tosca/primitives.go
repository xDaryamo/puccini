package tosca

import (
	"fmt"
	"time"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca/reflection"
)

var PrimitiveTypeValidators = map[string]reflection.TypeValidator{
	"boolean": reflection.IsBool,
	"integer": reflection.IsInt,
	"float":   reflection.IsFloat,
	"string":  reflection.IsString,
	"time":    reflection.IsTime, // must not conflict with TOSCA "timestamp" name
	"list":    reflection.IsSliceOfStruct,
	"map":     reflection.IsMap,
}

var PrimitiveTypeZeroes = map[string]interface{}{
	"boolean": false,
	"integer": 0,
	"float":   0.0,
	"string":  "",
	"time":    time.Time{},
	"list":    ard.List{},
	"map":     make(ard.Map),
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
		return "time"
	case ard.List:
		return "list"
	case ard.Map:
		return "map"
	default:
		return fmt.Sprintf("%T", value)
	}
}
