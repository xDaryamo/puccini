package format

import (
	"fmt"

	"github.com/op/go-logging"
)

//
// Log
//

type Log struct {
	logger *logging.Logger
	prefix string
}

func NewLog(logger *logging.Logger, name string) *Log {
	return &Log{
		logger: logger,
		prefix: fmt.Sprintf("{%s} ", name),
	}
}

func (self *Log) Errorf(fmt string, args ...interface{}) {
	self.logger.Errorf(self.prefix+fmt, args...)
}

func (self *Log) Warningf(fmt string, args ...interface{}) {
	self.logger.Warningf(self.prefix+fmt, args...)
}

func (self *Log) Infof(fmt string, args ...interface{}) {
	self.logger.Infof(self.prefix+fmt, args...)
}

func (self *Log) Debugf(fmt string, args ...interface{}) {
	self.logger.Debugf(self.prefix+fmt, args...)
}
