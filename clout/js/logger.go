package js

import (
	"fmt"

	"github.com/tliron/kutil/logging"
)

//
// Logger
//

type Logger struct {
	Prefix string

	logger logging.Logger
}

func NewLogger(logger logging.Logger, name string) *Logger {
	return &Logger{
		Prefix: fmt.Sprintf("{%s} ", name),
		logger: logger,
	}
}

func (self *Logger) Errorf(format string, args ...interface{}) {
	self.logger.Errorf("%s"+format, self.prefixedArgs(args)...)
}

func (self *Logger) Warningf(format string, args ...interface{}) {
	self.logger.Warningf("%s"+format, self.prefixedArgs(args)...)
}

func (self *Logger) Noticef(format string, args ...interface{}) {
	self.logger.Noticef("%s"+format, self.prefixedArgs(args)...)
}

func (self *Logger) Infof(format string, args ...interface{}) {
	self.logger.Infof("%s"+format, self.prefixedArgs(args)...)
}

func (self *Logger) Debugf(format string, args ...interface{}) {
	self.logger.Debugf("%s"+format, self.prefixedArgs(args)...)
}

func (self *Logger) prefixedArgs(args []interface{}) []interface{} {
	args = append(args, nil)
	copy(args[1:], args)
	args[0] = self.Prefix
	return args
}
