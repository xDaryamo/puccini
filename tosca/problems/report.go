package problems

import (
	"fmt"
	"strings"
)

func (self *Problems) ReportInSection(skip int, message string, section string) bool {
	// We want our reports to fit in one line
	message = strings.ReplaceAll(message, "\n", "¶")

	return self.Append(NewProblem(message, section, skip+1))
}

func (self *Problems) Report(skip int, message string) bool {
	return self.ReportInSection(skip+1, message, "")
}

func (self *Problems) Reportf(skip int, format string, arg ...interface{}) bool {
	return self.Report(skip+1, fmt.Sprintf(format, arg...))
}

func (self *Problems) ReportProblematic(skip int, problematic Problematic) bool {
	return self.ReportInSection(skip+1, problematic.ProblemMessage(), problematic.ProblemSection())
}

func (self *Problems) ReportError(err error) bool {
	if problematic, ok := err.(Problematic); ok {
		return self.ReportProblematic(1, problematic)
	} else {
		return self.Reportf(1, "%s", err.Error())
	}
}
