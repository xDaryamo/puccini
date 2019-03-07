package hot

import (
	"github.com/op/go-logging"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca"
)

var log = logging.MustGetLogger("grammars.hot")

var Grammar = make(tosca.Grammar)

var DefaultScriptNamespace = make(tosca.ScriptNamespace)

func init() {
	Grammar["ServiceTemplate"] = ReadTemplate

	Grammar["Condition"] = ReadCondition
	Grammar["ConditionDefinition"] = ReadConditionDefinition
	Grammar["Constraint"] = ReadConstraint
	Grammar["Data"] = ReadData
	Grammar["Output"] = ReadOutput
	Grammar["Parameter"] = ReadParameter
	Grammar["ParameterGroup"] = ReadParameterGroup
	Grammar["Resource"] = ReadResource
	Grammar["Template"] = ReadTemplate
	Grammar["Value"] = ReadValue

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
