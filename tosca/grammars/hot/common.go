package hot

import (
	"github.com/op/go-logging"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca"
)

var log = logging.MustGetLogger("grammars.hot")

var Grammar = make(tosca.Grammar)

var DefaultScriptletNamespace = make(tosca.ScriptletNamespace)

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

	for name, scriptlet := range FunctionSourceCode {
		DefaultScriptletNamespace[name] = &tosca.Scriptlet{
			Scriptlet: js.CleanupScriptlet(scriptlet),
		}
	}

	for name, scriptlet := range ConstraintSourceCode {
		DefaultScriptletNamespace[name] = &tosca.Scriptlet{
			Scriptlet: js.CleanupScriptlet(scriptlet),
		}
	}
}
