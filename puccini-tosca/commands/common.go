package commands

import (
	"github.com/op/go-logging"
	"github.com/tebeka/atexit"
	"github.com/tliron/puccini/common"
	formatpkg "github.com/tliron/puccini/common/format"
	problemspkg "github.com/tliron/puccini/common/problems"
	"github.com/tliron/puccini/common/terminal"
)

var log = logging.MustGetLogger("puccini-tosca")

func FailOnProblems(problems *problemspkg.Problems) {
	if !problems.Empty() {
		if !terminal.Quiet {
			if problemsFormat != "" {
				if strict {
					ard, err := problems.ARD()
					common.FailOnError(err)
					formatpkg.Print(ard, problemsFormat, terminal.Stderr, strict, pretty)
				} else {
					formatpkg.Print(problems, problemsFormat, terminal.Stderr, strict, pretty)
				}
			} else {
				problems.Print(verbose > 0)
			}
		}
		atexit.Exit(1)
	}
}
