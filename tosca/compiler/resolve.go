package compiler

import (
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/clout/js"
	problemspkg "github.com/tliron/puccini/tosca/problems"
)

func Resolve(clout *cloutpkg.Clout, problems *problemspkg.Problems, format string, strict bool, pretty bool) {
	context := js.NewContext("tosca.resolve", log, true, format, strict, pretty, "")
	if err := context.Exec(clout, "tosca.resolve", map[string]interface{}{"problems": problems}); err != nil {
		problems.ReportError(err)
	}
}
