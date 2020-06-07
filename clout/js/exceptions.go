package js

import (
	"fmt"

	"github.com/dop251/goja"
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
