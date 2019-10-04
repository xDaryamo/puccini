package js

import (
	"github.com/tliron/puccini/clout"
)

func (self *Context) Exec(clout_ *clout.Clout, name string, apis map[string]interface{}) error {
	scriptlet, err := GetScriptlet(name, clout_)
	if err != nil {
		return err
	}

	program, err := GetProgram(name, scriptlet)
	if err != nil {
		return err
	}

	_, runtime := self.NewCloutContext(clout_)
	for name, api := range apis {
		runtime.Set(name, api)
	}
	_, err = runtime.RunProgram(program)

	return UnwrapException(err)
}
