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

func (self *Problems) ReportWithURL(message string, url string) {
	self.Append(Problem{Message: message, URL: url})
}

func (self *Problems) Reportf(format string, arg ...interface{}) {
	self.Report(fmt.Sprintf(format, arg...))
}

func (self *Problems) ReportError(err error) {
	self.Reportf("%s", err)
}
