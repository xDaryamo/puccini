package js

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/tliron/puccini/common/terminal"
)

func UnwrapException(err error) error {
	if exception, ok := err.(*goja.Exception); ok {
		original := exception.Value().Export()
		if wrapped, ok := original.(error); ok {
			return wrapped
		} else if map_, ok := original.(map[string]interface{}); ok {
			if value, ok := map_["value"]; ok {
				if wrapped, ok := value.(error); ok {
					return wrapped
				} else {
					return fmt.Errorf("%s", value)
				}
			}
		} else {
			return fmt.Errorf("%s", original)
		}
	}
	return err
}

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
	r := fmt.Sprintf("%s: call to %s failed", self.FunctionCall.Path, self.Signature())
	if self.Message != "" {
		r += fmt.Sprintf(", %s", self.Message)
	}
	if self.Cause != nil {
		r += fmt.Sprintf(" because %s", self.Cause.Error())
	}
	return r
}

// fmt.Stringer interface
func (self *Error) String() string {
	return self.Error()
}

// tosca.problems.Problematic interface
func (self *Error) Problem() (string, string, int, int) {
	r := fmt.Sprintf("%s: call to %s failed", terminal.ColorPath(self.FunctionCall.Path), terminal.ColorName(self.Signature()))
	if self.Message != "" {
		r += fmt.Sprintf(", %s", self.Message)
	}
	if self.Cause != nil {
		if jsError, ok := self.Cause.(*Error); ok {
			message, _, _, _ := jsError.Problem()
			r += fmt.Sprintf(" because %s", message)
		} else {
			r += fmt.Sprintf(" because %s", self.Cause.Error())
		}
	}
	return r, self.FunctionCall.URL, self.FunctionCall.Row, self.FunctionCall.Column
}
