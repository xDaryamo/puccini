package js

import (
	"fmt"

	"github.com/dop251/goja"
)

func CallGojaFunction(runtime *goja.Runtime, function goja.Value, this any, arguments []any) (any, error) {
	function_, ok := goja.AssertFunction(function)
	if !ok {
		return nil, fmt.Errorf("not a function: %v", function)
	}

	values := make([]goja.Value, len(arguments))
	for index, argument := range arguments {
		values[index] = runtime.ToValue(argument)
	}

	r, err := function_(runtime.ToValue(this), values...)
	if err != nil {
		return nil, UnwrapException(err)
	}

	return r.Export(), nil
}

func UnwrapException(err error) error {
	if exception, ok := err.(*goja.Exception); ok {
		original := exception.Value().Export()
		if wrapped, ok := original.(error); ok {
			return wrapped
		} else if map_, ok := original.(map[string]any); ok {
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
