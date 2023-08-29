package commands

import (
	"os"

	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
	problemspkg "github.com/tliron/kutil/problems"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/transcribe"
	"github.com/tliron/kutil/util"
)

const toolName = "puccini-tosca"

var log = commonlog.GetLogger(toolName)

var importPaths []string
var template string
var inputs map[string]string
var inputsUrl string
var inputValues = make(map[string]any)
var problemsFormat string
var quirks []string
var urlMappings map[string]string

func FailOnProblems(problems *problemspkg.Problems) {
	if !problems.Empty() {
		if !terminal.Quiet {
			if problemsFormat != "" {
				transcribe.Print(problems, problemsFormat, os.Stderr, strict, pretty)
			} else {
				problems.Print(verbose > 0)
			}
		}
		util.Exit(1)
	}
}

func Bases(urlContext *exturl.Context, withImportPaths bool) []exturl.URL {
	var bases []exturl.URL

	if withImportPaths {
		for _, importPath := range importPaths {
			bases = append(bases, urlContext.NewAnyOrFileURL(importPath))
		}
	}

	workingDir, err := urlContext.NewWorkingDirFileURL()
	util.FailOnError(err)
	bases = append(bases, workingDir)

	return bases
}
