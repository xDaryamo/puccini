package cmd

import (
	"github.com/op/go-logging"
	"github.com/tebeka/atexit"
	"github.com/tliron/puccini/common/terminal"
	"github.com/tliron/puccini/tosca/problems"
)

var log = logging.MustGetLogger("puccini-tosca")

func FailOnProblems(problems_ *problems.Problems) {
	if !problems_.Empty() {
		if !terminal.Quiet {
			problems_.Print(verbose > 0)
		}
		atexit.Exit(1)
	}
}
