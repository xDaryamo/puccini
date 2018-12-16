package js

import (
	"encoding/json"
	"fmt"

	"github.com/dop251/goja"
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/clout"
)

//
// CloutContext
//

type CloutContext struct {
	*clout.Clout

	Context *Context
	Runtime *goja.Runtime
}

func (self *Context) NewCloutContext(c *clout.Clout) (*CloutContext, *goja.Runtime) {
	runtime := self.NewRuntime()
	context := &CloutContext{
		Clout:   c,
		Context: self,
		Runtime: runtime,
	}
	runtime.Set("clout", NewCloutApi(context))
	return context, runtime
}

func (self *CloutContext) CallFunction(site interface{}, source interface{}, target interface{}, name string, functionName string, arguments []interface{}) (interface{}, error) {
	sourceCode, err := GetScriptSourceCode(name, self.Clout)
	if err != nil {
		return nil, err
	}

	program, err := GetProgram(name, sourceCode)
	if err != nil {
		return nil, err
	}

	runtime := self.NewRuntime()
	runtime.Set("site", site)
	runtime.Set("source", source)
	runtime.Set("target", target)

	_, err = runtime.RunProgram(program)
	if err != nil {
		return nil, err
	}

	return CallFunction(runtime, functionName, arguments)
}

func (self *CloutContext) NewRuntime() *goja.Runtime {
	_, runtime := self.Context.NewCloutContext(self.Clout)
	return runtime
}

//
// CloutApi
//

type CloutApi struct {
	*clout.Clout

	context *CloutContext
}

func NewCloutApi(context *CloutContext) *CloutApi {
	return &CloutApi{context.Clout, context}
}

func (self *CloutApi) NewKey() string {
	return clout.NewKey()
}

func (self *CloutApi) Exec(name string) error {
	sourceCode, err := GetScriptSourceCode(name, self.context.Clout)
	if err != nil {
		return err
	}

	program, err := GetProgram(name, sourceCode)
	if err != nil {
		return err
	}

	_, err = self.context.Runtime.RunProgram(program)
	return err
}

func (self *CloutApi) NewCoercible(value goja.Value, site interface{}, source interface{}, target interface{}) (Coercible, error) {
	if goja.IsUndefined(value) {
		return nil, fmt.Errorf("undefined")
	}

	if coercible, err := self.context.NewCoercible(value.Export(), site, source, target); err == nil {
		return coercible, nil
	} else {
		return nil, err
	}
}

func (self *CloutApi) NewConstraints(value goja.Value, site interface{}, source interface{}, target interface{}) (Constraints, error) {
	if goja.IsUndefined(value) {
		return nil, fmt.Errorf("undefined")
	}

	exported := value.Export()
	if list_, ok := exported.(ard.List); ok {
		if constraints, err := self.context.NewConstraints(list_, site, source, target); err == nil {
			return constraints, nil
		} else {
			return nil, err
		}
	}

	return nil, fmt.Errorf("not an array")
}

func (self *CloutApi) Coerce(value interface{}) (interface{}, error) {
	if coercible, ok := value.(Coercible); ok {
		return coercible.Coerce()
	}

	return value, nil
}

func (self *CloutApi) Unwrap(value interface{}) interface{} {
	if coercible, ok := value.(Coercible); ok {
		return coercible.Unwrap()
	}

	return value
}

func (self *CloutApi) GetPlugins(name string) []goja.Value {
	plugins, err := GetPlugins(name, self.context)
	self.context.Context.FailOnError(err)
	return plugins
}

// json.Marshaler interface
func (self *CloutApi) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Clout)
}

// yaml.Marshaler interface
func (self *CloutApi) MarshalYAML() (interface{}, error) {
	return self.Clout, nil
}
