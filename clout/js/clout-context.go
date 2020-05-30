package js

import (
	"github.com/dop251/goja"
	cloutpkg "github.com/tliron/puccini/clout"
)

//
// CloutContext
//

type CloutContext struct {
	Context *Context
	Clout   *cloutpkg.Clout
	Runtime *goja.Runtime
}

func (self *Context) NewCloutContext(clout *cloutpkg.Clout, runtime *goja.Runtime) *CloutContext {
	return &CloutContext{
		Context: self,
		Clout:   clout,
		Runtime: runtime,
	}
}

func (self *CloutContext) Exec(scriptletName string) error {
	scriptlet, err := GetScriptlet(scriptletName, self.Clout)
	if err != nil {
		return err
	}

	program, err := self.Context.GetProgram(scriptletName, scriptlet)
	if err != nil {
		return err
	}

	_, err = self.Runtime.RunProgram(program)

	return UnwrapException(err)
}

func (self *CloutContext) ExecAll(scriptletBaseName string) error {
	if scriptletNames, err := GetScriptletNamesInSection(scriptletBaseName, self.Clout); err == nil {
		for _, scriptletName := range scriptletNames {
			if err := self.Exec(scriptletName); err != nil {
				return err
			}
		}
	}

	return nil
}

func (self *CloutContext) NewRuntime(apis map[string]interface{}) *goja.Runtime {
	return self.Context.NewCloutRuntime(self.Clout, apis)
}

func (self *CloutContext) CallFunction(scriptletName string, functionName string, arguments []interface{}, functionCallContext FunctionCallContext) (interface{}, error) {
	scriptlet, err := GetScriptlet(scriptletName, self.Clout)
	if err != nil {
		return nil, err
	}

	program, err := self.Context.GetProgram(scriptletName, scriptlet)
	if err != nil {
		return nil, err
	}

	runtime := self.NewRuntime(functionCallContext.API())

	_, err = runtime.RunProgram(program)
	if err != nil {
		return nil, UnwrapException(err)
	}

	return CallFunction(runtime, functionName, arguments)
}
