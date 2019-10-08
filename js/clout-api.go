package js

import (
	"encoding/json"
	"errors"

	"github.com/dop251/goja"
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/clout"
)

//
// CloutApi
//

type CloutApi struct {
	*clout.Clout

	cloutContext *CloutContext
}

func (self *Context) NewCloutApi(clout_ *clout.Clout, runtime *goja.Runtime) *CloutApi {
	return &CloutApi{
		clout_,
		self.NewCloutContext(clout_, runtime),
	}
}

func (self *CloutApi) NewKey() string {
	return clout.NewKey()
}

func (self *CloutApi) Exec(scriptletName string) error {
	return self.cloutContext.Exec(scriptletName)
}

func (self *CloutApi) NewCoercible(value goja.Value, site interface{}, source interface{}, target interface{}) (Coercible, error) {
	if goja.IsUndefined(value) {
		return nil, errors.New("undefined")
	}

	if coercible, err := self.cloutContext.NewCoercible(value.Export(), FunctionCallContext{site, source, target}); err == nil {
		return coercible, nil
	} else {
		return nil, err
	}
}

func (self *CloutApi) NewConstraints(value goja.Value, site interface{}, source interface{}, target interface{}) (Constraints, error) {
	if goja.IsUndefined(value) {
		return nil, errors.New("undefined")
	}

	exported := value.Export()
	if list_, ok := exported.(ard.List); ok {
		if constraints, err := self.cloutContext.NewConstraints(list_, FunctionCallContext{site, source, target}); err == nil {
			return constraints, nil
		} else {
			return nil, err
		}
	}

	return nil, errors.New("not an array")
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
	plugins, err := GetPlugins(name, self.cloutContext)
	self.cloutContext.Context.FailOnError(err)
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
