package hot

import (
	"github.com/op/go-logging"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca"
)

var log = logging.MustGetLogger("grammars.hot")

var Readers = make(map[string]tosca.Reader)

var DefaultScriptNamespace = make(tosca.ScriptNamespace)

func init() {
	Readers["Condition"] = ReadCondition
	Readers["ConditionDefinition"] = ReadConditionDefinition
	Readers["Constraint"] = ReadConstraint
	Readers["Data"] = ReadData
	Readers["Output"] = ReadOutput
	Readers["Parameter"] = ReadParameter
	Readers["ParameterGroup"] = ReadParameterGroup
	Readers["Resource"] = ReadResource
	Readers["Template"] = ReadTemplate
	Readers["Value"] = ReadValue

	for name, sourceCode := range FunctionSourceCode {
		DefaultScriptNamespace[name] = &tosca.Script{
			SourceCode: js.Cleanup(sourceCode),
		}
	}

	for name, sourceCode := range ConstraintSourceCode {
		DefaultScriptNamespace[name] = &tosca.Script{
			SourceCode: js.Cleanup(sourceCode),
		}
	}
}
