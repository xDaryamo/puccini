package js

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/tliron/puccini/format"
)

func UnwrapException(err error) error {
	if exception, ok := err.(*goja.Exception); ok {
		original := exception.Value().Export()
		if map_, ok := original.(map[string]interface{}); ok {
			if value, ok := map_["value"]; ok {
				if wrapped, ok := value.(error); ok {
					return wrapped
				} else {
					return fmt.Errorf("%s", value)
				}
			}
		} else {
			if wrapped, ok := original.(error); ok {
				return wrapped
			} else {
				return fmt.Errorf("%s", original)
			}
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
func (self Error) Error() string {
	r := fmt.Sprintf("%s: call to %s failed", self.FunctionCall.Path, self.Signature())
	if self.Message != "" {
		r += fmt.Sprintf(", %s", self.Message)
	}
	if self.Cause != nil {
		r += fmt.Sprintf(" due to %s", self.Cause.Error())
	}
	return r
}

// tosca.problems.Problematic interface
func (self Error) ProblemMessage() string {
	r := fmt.Sprintf("%s: call to %s failed", format.ColorPath(self.FunctionCall.Path), format.ColorName(self.Signature()))
	if self.Message != "" {
		r += fmt.Sprintf(", %s", self.Message)
	}
	if self.FunctionCall.Location != "" {
		if r != "" {
			r += " "
		}
		r += format.ColorValue("@" + self.FunctionCall.Location)
	}
	if self.Cause != nil {
		if jsError, ok := self.Cause.(*Error); ok {
			r += fmt.Sprintf(" due to %s", jsError.ProblemMessage())
		} else {
			r += fmt.Sprintf(" due to %s", self.Cause.Error())
		}
	}
	return r
}

// tosca.problems.Problematic interface
func (self Error) ProblemSection() string {
	return self.FunctionCall.URL
}
