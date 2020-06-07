package js

import (
	"fmt"

	"github.com/tliron/puccini/common/terminal"
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
	r := fmt.Sprintf("%s: ", self.FunctionCall.Path)
	if self.Message != "" {
		r += fmt.Sprintf("%s in call to %s", self.Message, self.Signature())
	} else {
		r += fmt.Sprintf("call to %s failed", self.Signature())
	}
	if self.Cause != nil {
		if jsError, ok := self.Cause.(*Error); ok {
			message, _, _, _ := jsError.Problem()
			r += fmt.Sprintf(" because %s", message)
		} else {
			r += fmt.Sprintf(" because %s", self.Cause.Error())
		}
	}
	return r
}

// fmt.Stringer interface
func (self *Error) String() string {
	return self.Error()
}

// problems.Problematic interface
func (self *Error) Problem() (string, string, int, int) {
	r := fmt.Sprintf("%s: ", terminal.ColorPath(self.FunctionCall.Path))
	if self.Message != "" {
		r += fmt.Sprintf("%s in call to %s", self.Message, terminal.ColorName(self.Signature()))
	} else {
		r += fmt.Sprintf("call to %s failed", terminal.ColorName(self.Signature()))
	}
	if self.Cause != nil {
		if jsError, ok := self.Cause.(*Error); ok {
			message, _, _, _ := jsError.Problem()
			r += fmt.Sprintf(" because %s", terminal.ColorError(message))
		} else {
			r += fmt.Sprintf(" because %s", terminal.ColorError(self.Cause.Error()))
		}
	}
	return r, self.FunctionCall.URL, self.FunctionCall.Row, self.FunctionCall.Column
}
