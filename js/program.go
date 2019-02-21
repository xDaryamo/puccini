package js

import (
	"sync"

	"github.com/dop251/goja"
)

func GetProgram(name string, script string) (*goja.Program, error) {
	p, ok := ProgramCache.Load(script)
	if !ok {
		program, err := goja.Compile(name, script, true)
		if err != nil {
			return nil, err
		}
		p, _ = ProgramCache.LoadOrStore(script, program)
	}

	return p.(*goja.Program), nil
}

var ProgramCache sync.Map
