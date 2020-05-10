package compiler

import (
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/clout/js"
	problemspkg "github.com/tliron/puccini/common/problems"
	urlpkg "github.com/tliron/puccini/url"
)

func Coerce(clout *cloutpkg.Clout, problems *problemspkg.Problems, urlContext *urlpkg.Context, format string, strict bool, allowTimestamps bool, pretty bool) {
	context := js.NewContext("tosca.coerce", log, true, format, strict, allowTimestamps, pretty, "", urlContext)
	if err := context.Exec(clout, "tosca.coerce", map[string]interface{}{"problems": problems}); err != nil {
		problems.ReportError(err)
	}
}
