package commands

import (
	"os"

	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
	"github.com/tliron/go-transcribe"
	problemspkg "github.com/tliron/kutil/problems"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

const toolName = "puccini-tosca"

var (
	log = commonlog.GetLogger(toolName)

	importPaths    []string
	template       string
	inputs         map[string]string
	inputsUrl      string
	inputValues    = make(map[string]any)
	problemsFormat string
	quirks         []string
	urlMappings    map[string]string
)

func Transcriber() *transcribe.Transcriber {
	return &transcribe.Transcriber{
		File:        output,
		Format:      format,
		ForTerminal: pretty,
		Strict:      strict,
		Base64:      base64,
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

func FailOnProblems(problems *problemspkg.Problems) {
	if !problems.Empty() {
		if !terminal.Quiet {
			if problemsFormat != "" {
				transcriber := Transcriber().Clone()
				transcriber.Writer = os.Stderr
				transcriber.Format = problemsFormat

				transcriber.Write(problems)
			} else {
				problems.Print(verbose > 0)
			}
		}
		util.Exit(1)
	}
}
