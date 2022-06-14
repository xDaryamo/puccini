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
	Exec("tosca.resolve", arguments, clout, problems, urlContext, format, strict, pretty)
}

func Coerce(clout *cloutpkg.Clout, problems *problemspkg.Problems, urlContext *urlpkg.Context, history bool, format string, strict bool, pretty bool) {
	var arguments map[string]string
	if !history {
		arguments = make(map[string]string)
		arguments["history"] = "false"
	}
	Exec("tosca.coerce", arguments, clout, problems, urlContext, format, strict, pretty)
}

func Outputs(clout *cloutpkg.Clout, problems *problemspkg.Problems, urlContext *urlpkg.Context, format string, strict bool, pretty bool) {
	Exec("tosca.outputs", nil, clout, problems, urlContext, format, strict, pretty)
}

func Exec(scriptletName string, arguments map[string]string, clout *cloutpkg.Clout, problems *problemspkg.Problems, urlContext *urlpkg.Context, format string, strict bool, pretty bool) {
	context := NewContext(scriptletName, log, arguments, true, format, strict, pretty, "", urlContext)
	if _, err := context.Require(clout, scriptletName, map[string]any{"problems": problems}); err != nil {
		problems.ReportError(err)
	}
}
