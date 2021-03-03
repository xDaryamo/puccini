package js

import (
	"fmt"

	"github.com/tliron/kutil/terminal"
)

//
// Error
//

type Error struct {
	FunctionCall *FunctionCall
	Arguments    []interface{}
	Message      string
	Cause        error
}

func (self *FunctionCall) NewError(arguments []interface{}, message string, cause error) *Error {
	return &Error{
		FunctionCall: self,
		Arguments:    arguments,
		Message:      message,
		Cause:        cause,
	}
}

func (self *FunctionCall) NewErrorf(arguments []interface{}, format string, arg ...interface{}) *Error {
	return self.NewError(arguments, fmt.Sprintf(format, arg...), nil)
}

func (self *FunctionCall) WrapError(arguments []interface{}, err error) *Error {
	return self.NewError(arguments, "", UnwrapException(err))
}

func (self *Error) Signature() string {
	return self.FunctionCall.Signature(self.Arguments)
}

// error interface
func (self *Error) Error() string {
	message := fmt.Sprintf("%s: ", self.FunctionCall.Path)
	if self.Message != "" {
		message += fmt.Sprintf("%s in call to %s", self.Message, self.Signature())
	} else {
		message += fmt.Sprintf("call to %s failed", self.Signature())
	}
	if self.Cause != nil {
		if jsError, ok := self.Cause.(*Error); ok {
			_, _, message_, _, _ := jsError.Problem()
			message += fmt.Sprintf(" because %s", message_)
		} else {
			message += fmt.Sprintf(" because %s", self.Cause.Error())
		}
	}
	return message
}

// fmt.Stringer interface
func (self *Error) String() string {
	return self.Error()
}

// problems.Problematic interface
func (self *Error) Problem() (string, string, string, int, int) {
	message := fmt.Sprintf("%s: ", terminal.StylePath(self.FunctionCall.Path))
	if self.Message != "" {
		message += fmt.Sprintf("%s in call to %s", self.Message, terminal.StyleName(self.Signature()))
	} else {
		message += fmt.Sprintf("call to %s failed", terminal.StyleName(self.Signature()))
	}
	if self.Cause != nil {
		if jsError, ok := self.Cause.(*Error); ok {
			_, _, message_, _, _ := jsError.Problem()
			message += fmt.Sprintf(" because %s", terminal.StyleError(message_))
		} else {
			message += fmt.Sprintf(" because %s", terminal.StyleError(self.Cause.Error()))
		}
	}
	return self.FunctionCall.URL, "", message, self.FunctionCall.Row, self.FunctionCall.Column
}
