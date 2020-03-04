package compiler

import (
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/clout/js"
	problemspkg "github.com/tliron/puccini/tosca/problems"
)

func Coerce(clout *cloutpkg.Clout, problems *problemspkg.Problems, format string, strict bool, pretty bool) {
	context := js.NewContext("tosca.coerce", log, true, format, strict, pretty, "")
	if err := context.Exec(clout, "tosca.coerce", map[string]interface{}{"problems": problems}); err != nil {
		problems.ReportError(err)
	}
}
