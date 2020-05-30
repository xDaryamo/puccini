package js

import (
	"fmt"

	"github.com/op/go-logging"
)

//
// Log
//

type Log struct {
	Prefix string

	logger *logging.Logger
}

func NewLog(logger *logging.Logger, name string) *Log {
	return &Log{
		Prefix: fmt.Sprintf("{%s} ", name),
		logger: logger,
	}
}

func (self *Log) Errorf(format string, args ...interface{}) {
	self.logger.Errorf("%s"+format, self.prefixedArgs(args)...)
}

func (self *Log) Warningf(format string, args ...interface{}) {
	self.logger.Warningf("%s"+format, self.prefixedArgs(args)...)
}

func (self *Log) Noticef(format string, args ...interface{}) {
	self.logger.Noticef("%s"+format, self.prefixedArgs(args)...)
}

func (self *Log) Infof(format string, args ...interface{}) {
	self.logger.Infof("%s"+format, self.prefixedArgs(args)...)
}

func (self *Log) Debugf(format string, args ...interface{}) {
	self.logger.Debugf("%s"+format, self.prefixedArgs(args)...)
}

func (self *Log) prefixedArgs(args []interface{}) []interface{} {
	args = append(args, nil)
	copy(args[1:], args)
	args[0] = self.Prefix
	return args
}
