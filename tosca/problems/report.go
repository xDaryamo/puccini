package problems

import (
	"fmt"
)

func (self *Problems) Append(problem Problem) {
	self.Problems = append(self.Problems, problem)
}

func (self *Problems) Report(message string) {
	self.Append(Problem{Message: message})
}

func (self *Problems) ReportInSection(message string, section string) {
	self.Append(Problem{Message: message, Section: section})
}

func (self *Problems) Reportf(format string, arg ...interface{}) {
	self.Report(fmt.Sprintf(format, arg...))
}

func (self *Problems) ReportError(err error) {
	if problematic, ok := err.(Problematic); ok {
		self.ReportProblematic(problematic)
	} else {
		self.Reportf("%s", err.Error())
	}
}

func (self *Problems) ReportProblematic(problematic Problematic) {
	self.ReportInSection(problematic.ProblemMessage(), problematic.ProblemSection())
}
