package js

import (
	"fmt"

	"github.com/dop251/goja"
)

func CallFunction(runtime *goja.Runtime, functionName string, arguments []interface{}) (interface{}, error) {
	value := runtime.Get(functionName)
	if value == nil {
		return nil, fmt.Errorf("scriptlet does not have a %q function", functionName)
	}

	function, ok := goja.AssertFunction(value)
	if !ok {
		return nil, fmt.Errorf("scriptlet has a %q variable but it's not a function", functionName)
	}

	values := make([]goja.Value, len(arguments))
	for index, argument := range arguments {
		values[index] = runtime.ToValue(argument)
	}

	r, err := function(nil, values...)
	if err != nil {
		return nil, UnwrapException(err)
	}

	return r.Export(), nil
}
