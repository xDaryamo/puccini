package hot

import (
	"github.com/tliron/kutil/logging"
	"github.com/tliron/puccini/tosca"
)

var log = logging.GetLogger("puccini.grammars.hot")
var logRender = logging.NewSubLogger(log, "render")
var logNormalize = logging.NewSubLogger(log, "normalize")

var Grammar = tosca.NewGrammar()

var DefaultScriptletNamespace = tosca.NewScriptletNamespace()

func init() {
	Grammar.RegisterVersion("heat_template_version", "train", "") // not mentioned in spec, but probably supported
	Grammar.RegisterVersion("heat_template_version", "stein", "") // not mentioned in spec, but probably supported
	Grammar.RegisterVersion("heat_template_version", "rocky", "")
	Grammar.RegisterVersion("heat_template_version", "queens", "")
	Grammar.RegisterVersion("heat_template_version", "pike", "")
	Grammar.RegisterVersion("heat_template_version", "newton", "")
	Grammar.RegisterVersion("heat_template_version", "ocata", "")
	Grammar.RegisterVersion("heat_template_version", "2018-08-31", "") // train, stein, rocky
	Grammar.RegisterVersion("heat_template_version", "2018-03-02", "") // queens
	Grammar.RegisterVersion("heat_template_version", "2017-09-01", "") // pike
	Grammar.RegisterVersion("heat_template_version", "2017-02-24", "") // ocata
	Grammar.RegisterVersion("heat_template_version", "2016-10-14", "") // newton
	Grammar.RegisterVersion("heat_template_version", "2016-04-08", "") // mitaka
	Grammar.RegisterVersion("heat_template_version", "2015-10-15", "") // liberty
	Grammar.RegisterVersion("heat_template_version", "2015-04-30", "") // kilo
	Grammar.RegisterVersion("heat_template_version", "2014-10-16", "") // juno
	Grammar.RegisterVersion("heat_template_version", "2013-05-23", "") // icehouse

	Grammar.RegisterReader("$Root", ReadTemplate)

	Grammar.RegisterReader("Condition", ReadCondition)
	Grammar.RegisterReader("ConditionDefinition", ReadConditionDefinition)
	Grammar.RegisterReader("Constraint", ReadConstraint)
	Grammar.RegisterReader("Data", ReadData)
	Grammar.RegisterReader("Output", ReadOutput)
	Grammar.RegisterReader("Parameter", ReadParameter)
	Grammar.RegisterReader("ParameterGroup", ReadParameterGroup)
	Grammar.RegisterReader("Resource", ReadResource)
	Grammar.RegisterReader("Template", ReadTemplate)
	Grammar.RegisterReader("Value", ReadValue)

	DefaultScriptletNamespace.RegisterScriptlets(FunctionScriptlets, nil)
	DefaultScriptletNamespace.RegisterScriptlets(ConstraintScriptlets, ConstraintNativeArgumentIndexes)
}
