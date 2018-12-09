package compiler

import (
	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca/problems"
)

func Coerce(clout_ *clout.Clout, problems_ *problems.Problems) *clout.Clout {
	context := js.NewContext("tosca.coerce", log, false, "yaml", "")
	err := context.Exec(clout_, "tosca.coerce", map[string]interface{}{"problems": problems_})
	if err != nil {
		problems_.ReportError(err)
	}

	return clout_
}
