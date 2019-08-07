package js

import (
	"github.com/tliron/puccini/clout"
)

func (self *Context) Exec(clout_ *clout.Clout, name string, apis map[string]interface{}) error {
	sourceCode, err := GetScriptSourceCode(name, clout_)
	if err != nil {
		return err
	}

	program, err := GetProgram(name, sourceCode)
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
