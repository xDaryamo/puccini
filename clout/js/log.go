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

func (self *Log) Errorf(f string, args ...interface{}) {
	self.logger.Errorf(self.Prefix+f, args...)
}

func (self *Log) Warningf(f string, args ...interface{}) {
	self.logger.Warningf(self.Prefix+f, args...)
}

func (self *Log) Infof(f string, args ...interface{}) {
	self.logger.Infof(self.Prefix+f, args...)
}

func (self *Log) Debugf(f string, args ...interface{}) {
	self.logger.Debugf(self.Prefix+f, args...)
}
