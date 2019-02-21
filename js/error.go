package js

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/tliron/puccini/format"
)

func UnwrapError(err error) error {
	if exception, ok := err.(*goja.Exception); ok {
		original := exception.Value().Export()
		if map_, ok := original.(map[string]interface{}); ok {
			if value, ok := map_["value"]; ok {
				return fmt.Errorf("%s", value)
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
	Function  *Function
	Arguments []interface{}
	Message   string
}

func (self *Function) NewError(arguments []interface{}, message string) *Error {
	return &Error{
		Function:  self,
		Arguments: arguments,
		Message:   message,
	}
}

func (self *Function) NewErrorf(arguments []interface{}, format string, arg ...interface{}) *Error {
	return self.NewError(arguments, fmt.Sprintf(format, arg...))
}

func (self *Function) WrapError(arguments []interface{}, err error) *Error {
	return self.NewError(arguments, UnwrapError(err).Error())
}

func (self *Error) Signature() string {
	return self.Function.Signature(self.Arguments)
}

// error interface
func (self Error) Error() string {
	if self.Message == "" {
		return fmt.Sprintf("%s: call to \"%s\" failed", self.Function.Path, self.Signature())
	} else {
		return fmt.Sprintf("%s: call to \"%s\" failed: %s", self.Function.Path, self.Signature(), self.Message)
	}
}

func (self Error) ColorError() string {
	if self.Message == "" {
		return fmt.Sprintf("%s: call to \"%s\" failed", format.ColorPath(self.Function.Path), format.ColorName(self.Signature()))
	} else {
		return fmt.Sprintf("%s: call to \"%s\" failed: %s", format.ColorPath(self.Function.Path), format.ColorName(self.Signature()), self.Message)
	}
}
