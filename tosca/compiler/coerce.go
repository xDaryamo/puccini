package compiler

import (
	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca/problems"
)

func Coerce(clout_ *clout.Clout, problems_ *problems.Problems, ardFormat string, pretty bool) {
	context := js.NewContext("tosca.coerce", log, true, ardFormat, pretty, "")
	if err := context.Exec(clout_, "tosca.coerce", map[string]interface{}{"problems": problems_}); err != nil {
		problems_.ReportError(err)
	}
}
