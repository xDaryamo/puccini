package problems

import (
	"fmt"
)

func (self *Problems) Report(message string) {
	*self = append(*self, Problem{Message: message})
}

func (self *Problems) Reportf(format string, arg ...interface{}) {
	self.Report(fmt.Sprintf(format, arg...))
}

func (self *Problems) ReportError(err error) {
	self.Reportf("%s", err)
}
