package compiler

import (
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/clout/js"
	problemspkg "github.com/tliron/puccini/common/problems"
)

func Coerce(clout *cloutpkg.Clout, problems *problemspkg.Problems, format string, strict bool, allowTimestamps bool, pretty bool) {
	context := js.NewContext("tosca.coerce", log, true, format, strict, allowTimestamps, pretty, "")
	if err := context.Exec(clout, "tosca.coerce", map[string]interface{}{"problems": problems}); err != nil {
		problems.ReportError(err)
	}
}
