package js

import (
	"errors"
	"fmt"
	"sort"

	"github.com/dop251/goja"
)

func GetPlugins(name string, cloutContext *CloutContext) ([]goja.Value, error) {
	scripts, err := GetScriptlets(name, cloutContext.Clout)
	if err != nil {
		return nil, nil
	}

	var plugins []goja.Value

	for _, value := range scripts {
		if _, ok := value.(string); !ok {
			return nil, fmt.Errorf("plugin script is not a string: %T", value)
		}
	}

	sort.Slice(scripts, func(i, j int) bool {
		return scripts[i].(string) < scripts[j].(string)
	})

	for _, value := range scripts {
		scriptlet := value.(string)

		program, err := cloutContext.Context.GetProgram("<plugin>", scriptlet)
		if err != nil {
			return nil, err
		}

		runtime := cloutContext.NewRuntime(nil)
		_, err = runtime.RunProgram(program)
		if err != nil {
			return nil, err
		}

		plugin := runtime.Get("plugin")
		if plugin == nil {
			return nil, errors.New("plugin script does not define \"plugin\" variable")
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
