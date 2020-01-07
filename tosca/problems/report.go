package problems

import (
	"fmt"
	"strings"
)

func (self *Problems) ReportInSection(message string, section string) bool {
	// We want our reports to fit in one line
	message = strings.ReplaceAll(message, "\n", "¶")

	return self.Append(Problem{Message: message, Section: section})
}

func (self *Problems) Report(message string) bool {
	return self.ReportInSection(message, "")
}

func (self *Problems) Reportf(format string, arg ...interface{}) bool {
	return self.Report(fmt.Sprintf(format, arg...))
}

func (self *Problems) ReportError(err error) bool {
	if problematic, ok := err.(Problematic); ok {
		return self.ReportProblematic(problematic)
	} else {
		return self.Reportf("%s", err.Error())
	}
}

func (self *Problems) ReportProblematic(problematic Problematic) bool {
	return self.ReportInSection(problematic.ProblemMessage(), problematic.ProblemSection())
}
