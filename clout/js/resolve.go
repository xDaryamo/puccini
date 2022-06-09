package js

import (
	problemspkg "github.com/tliron/kutil/problems"
	urlpkg "github.com/tliron/kutil/url"
	cloutpkg "github.com/tliron/puccini/clout"
)

func Resolve(clout *cloutpkg.Clout, problems *problemspkg.Problems, urlContext *urlpkg.Context, history bool, format string, strict bool, pretty bool) {
	var arguments map[string]string
	if !history {
		arguments = make(map[string]string)
		arguments["history"] = "false"
	}
	context := NewContext("tosca.resolve", log, arguments, true, format, strict, pretty, "", urlContext)
	if _, err := context.Require(clout, "tosca.resolve", map[string]any{"problems": problems}); err != nil {
		problems.ReportError(err)
	}
}
