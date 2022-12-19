package js

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/logging"
)

var log = logging.GetLogger("puccini.js")
var logEvaluate = logging.NewScopeLogger(log, "evaluate")
var logValidate = logging.NewScopeLogger(log, "validate")
var logConvert = logging.NewScopeLogger(log, "convert")

func asInt(value any) (int, bool) {
	switch value_ := value.(type) {
	case int64:
		return int(value_), true
	case int32:
		return int(value_), true
	case int16:
		return int(value_), true
	case int8:
		return int(value_), true
	case int:
		return value_, true
	case uint64:
		return int(value_), true
	case uint32:
		return int(value_), true
	case uint16:
		return int(value_), true
	case uint8:
		return int(value_), true
	case uint:
		return int(value_), true
	case float64:
		return int(value_), true
	case float32:
		return int(value_), true
	}
	return 0, false
}

func asList(value any) (ard.List, bool) {
	if list, ok := value.(ard.List); ok {
		return list, true
	} else if value == nil {
		return nil, true
	} else {
		return nil, false
	}
}

func asStringMap(value any) (ard.StringMap, bool) {
	if map_, ok := value.(ard.StringMap); ok {
		return map_, true
	} else if value == nil {
		return make(ard.StringMap), true
	} else {
		return nil, false
	}
}
