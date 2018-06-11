package js

import (
	"encoding/json"
	"fmt"

	"github.com/dop251/goja"
	"github.com/tliron/puccini/clout"
)

//
// Clout
//

type Clout struct {
	*clout.Clout
	Runtime *goja.Runtime
}

func NewClout(c *clout.Clout, runtime *goja.Runtime) *Clout {
	return &Clout{
		Clout:   c,
		Runtime: runtime,
	}
}

func (self *Clout) Exec(name string) error {
	sourceCode, err := GetScriptSourceCode(name, self.Clout)
	if err != nil {
		return err
	}

	program, err := GetProgram(name, sourceCode)
	if err != nil {
		return err
	}

	_, err = self.Runtime.RunProgram(program)
	return err
}

func (self *Clout) Prepare(value goja.Value, site interface{}, source interface{}, target interface{}) (Coercible, error) {
	if goja.IsUndefined(value) {
		return nil, fmt.Errorf("undefined")
	}
	coercible, err := NewCoercible(value.Export(), site, source, target, self.Clout)
	if err != nil {
		return nil, err
	}
	return coercible, nil
}

func (self *Clout) Coerce(value interface{}) (interface{}, error) {
	coercible, ok := value.(Coercible)
	if !ok {
		return value, nil
	}
	return coercible.Coerce()
}

// json.Marshaler interface
func (self *Clout) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Clout)
}

// yaml.Marshaler interface
func (self *Clout) MarshalYAML() (interface{}, error) {
	return self.Clout, nil
}
