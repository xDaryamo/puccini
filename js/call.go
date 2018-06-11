package js

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/tliron/puccini/clout"
)

func Call(runtime *goja.Runtime, functionName string, arguments []interface{}) (interface{}, error) {
	value := runtime.Get(functionName)
	if value == nil {
		return nil, fmt.Errorf("script does not have a \"%s\" function", functionName)
	}

	function, ok := goja.AssertFunction(value)
	if !ok {
		return nil, fmt.Errorf("script has a \"%s\" variable but it's not a function", functionName)
	}

	values := make([]goja.Value, 0, len(arguments))
	for _, argument := range arguments {
		values = append(values, runtime.ToValue(argument))
	}

	r, err := function(nil, values...)
	if err != nil {
		return nil, err
	}

	return r.Export(), nil
}

func CallClout(c *clout.Clout, site interface{}, source interface{}, target interface{}, name string, functionName string, arguments []interface{}) (interface{}, error) {
	sourceCode, err := GetScriptSourceCode(name, c)
	if err != nil {
		return nil, err
	}

	program, err := GetProgram(name, sourceCode)
	if err != nil {
		return nil, err
	}

	runtime := NewCloutRuntime(name, c)
	runtime.Set("site", site)
	runtime.Set("source", source)
	runtime.Set("target", target)

	_, err = runtime.RunProgram(program)
	if err != nil {
		return nil, err
	}

	return Call(runtime, functionName, arguments)
}
