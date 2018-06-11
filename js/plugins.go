package js

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/tliron/puccini/clout"
)

func GetPlugins(name string, c *clout.Clout) ([]goja.Value, error) {
	scripts, err := GetScripts(name, c)
	if err != nil {
		return nil, nil
	}

	var plugins []goja.Value

	for _, value := range scripts {
		sourceCode, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("plugin script is not a string: %T", value)
		}

		program, err := GetProgram("<plugin>", sourceCode)
		if err != nil {
			return nil, err
		}

		runtime := NewRuntime("plugin")
		_, err = runtime.RunProgram(program)
		if err != nil {
			return nil, err
		}

		plugin := runtime.Get("plugin")
		if plugin == nil {
			return nil, fmt.Errorf("plugin script does not define \"plugin\" variable")
		}

		plugins = append(plugins, plugin)
	}

	return plugins, nil

	//`x = function() {
	//	printf('hi from JS!!\n');
	//	return 3;
	//};`
	//
	// x := vm.Get("x")
	// y, _ := goja.AssertFunction(x)
	// z, _ := y(nil)
	// fmt.Printf("%s\n", z)
}
