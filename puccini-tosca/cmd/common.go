package cmd

import (
	"github.com/op/go-logging"
	"github.com/tebeka/atexit"
	format_ "github.com/tliron/puccini/common/format"
	"github.com/tliron/puccini/common/terminal"
	problems_ "github.com/tliron/puccini/tosca/problems"
)

var log = logging.MustGetLogger("puccini-tosca")

func FailOnProblems(problems *problems_.Problems) {
	if !problems.Empty() {
		if !terminal.Quiet {
			if problemsFormat != "" {
				format_.Print(problems, problemsFormat, terminal.Stderr, pretty)
			} else {
				problems.Print(verbose > 0)
			}
		}
		atexit.Exit(1)
	}
}
