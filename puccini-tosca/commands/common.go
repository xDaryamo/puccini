package commands

import (
	formatpkg "github.com/tliron/kutil/format"
	"github.com/tliron/kutil/logging"
	problemspkg "github.com/tliron/kutil/problems"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
)

const toolName = "puccini-tosca"

var log = logging.GetLogger(toolName)

func FailOnProblems(problems *problemspkg.Problems) {
	if !problems.Empty() {
		if !terminal.Quiet {
			if problemsFormat != "" {
				formatpkg.Print(problems, problemsFormat, terminal.Stderr, strict, pretty)
			} else {
				problems.Print(verbose > 0)
			}
		}
		util.Exit(1)
	}
}
