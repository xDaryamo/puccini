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

func (self *Log) Errorf(f string, args ...interface{}) {
	self.logger.Errorf(self.prefix+f, args...)
}

func (self *Log) Warningf(f string, args ...interface{}) {
	self.logger.Warningf(self.prefix+f, args...)
}

func (self *Log) Infof(f string, args ...interface{}) {
	self.logger.Infof(self.prefix+f, args...)
}

func (self *Log) Debugf(f string, args ...interface{}) {
	self.logger.Debugf(self.prefix+f, args...)
}
