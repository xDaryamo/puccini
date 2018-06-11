package v1_1

import (
	"github.com/op/go-logging"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca"
)

var log = logging.MustGetLogger("grammars.v1_1")

var Readers = make(map[string]tosca.Reader)

var DefaultScriptNamespace = make(tosca.ScriptNamespace)

func init() {
	for name, sourceCode := range FunctionSourceCode {
		DefaultScriptNamespace[name] = &tosca.Script{
			SourceCode: js.Cleanup(sourceCode),
		}
	}

	for name, sourceCode := range ConstraintClauseSourceCode {
		nativeArgumentIndexes, _ := ConstraintClauseNativeArgumentIndexes[name]
		DefaultScriptNamespace[name] = &tosca.Script{
			SourceCode:            js.Cleanup(sourceCode),
			NativeArgumentIndexes: nativeArgumentIndexes,
		}
	}
}

func CompareUint32(v1 uint32, v2 uint32) int {
	if v1 < v2 {
		return -1
	} else if v2 > v1 {
		return 1
	}
	return 0
}

func CompareUint64(v1 uint64, v2 uint64) int {
	if v1 < v2 {
		return -1
	} else if v2 > v1 {
		return 1
	}
	return 0
}

func CompareInt64(v1 int64, v2 int64) int {
	if v1 < v2 {
		return -1
	} else if v2 > v1 {
		return 1
	}
	return 0
}

func CompareFloat64(v1 float64, v2 float64) int {
	if v1 < v2 {
		return -1
	} else if v2 > v1 {
		return 1
	}
	return 0
}
