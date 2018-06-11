package tosca

import (
	"fmt"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca/reflection"
)

var PrimitiveTypeValidators = map[string]reflection.TypeValidator{
	"boolean": reflection.IsBool,
	"integer": reflection.IsInt,
	"float":   reflection.IsFloat,
	"string":  reflection.IsString,
	"map":     reflection.IsMap,
	"list":    reflection.IsSliceOfStruct,
}

var PrimitiveTypeZeroes = map[string]interface{}{
	"boolean": false,
	"integer": 0,
	"float":   0.0,
	"string":  "",
	"map":     make(ard.Map),
	"list":    ard.List{},
}

func PrimitiveTypeName(value interface{}) string {
	switch value.(type) {
	case bool:
		return "boolean"
	case int:
		return "integer"
	case float64: // YAML parser returns float64
		return "float"
	case ard.List:
		return "list"
	case ard.Map:
		return "map"
	default:
		return fmt.Sprintf("%T", value)
	}
}
