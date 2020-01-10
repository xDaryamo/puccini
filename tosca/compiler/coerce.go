package compiler

import (
	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/clout/js"
	"github.com/tliron/puccini/tosca/problems"
)

func Coerce(clout_ *clout.Clout, problems_ *problems.Problems, format string, pretty bool) {
	context := js.NewContext("tosca.coerce", log, true, format, pretty, "")
	if err := context.Exec(clout_, "tosca.coerce", map[string]interface{}{"problems": problems_}); err != nil {
		problems_.ReportError(err)
	}
}
