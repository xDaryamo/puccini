package js

import (
	"fmt"
	"strings"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/format"
)

//
// Function
//

type Function struct {
	Context     *CloutContext `json:"-" yaml:"-"`
	Name        string        `json:"function" yaml:"function"`
	Path        string        `json:"path" yaml:"path"`
	Arguments   []Coercible   `json:"arguments" yaml:"arguments"`
	Constraints Constraints   `json:"constraints" yaml:"constraints"`
	Site        interface{}   `json:"-" yaml:"-"`
	Source      interface{}   `json:"-" yaml:"-"`
	Target      interface{}   `json:"-" yaml:"-"`
	Notation    ard.Map       `json:"-" yaml:"-"`
}

func (self *CloutContext) NewFunction(data interface{}, site interface{}, source interface{}, target interface{}) (*Function, error) {
	map_, ok := data.(ard.Map)
	if !ok {
		return nil, fmt.Errorf("not a function")
	}

	function, ok := map_["function"]
	if !ok {
		return nil, fmt.Errorf("not a function")
	}

	c := Function{
		Context:  self,
		Site:     site,
		Source:   source,
		Target:   target,
		Notation: map_,
	}

	f, ok := function.(ard.Map)
	if !ok {
		return nil, fmt.Errorf("malformed function: not a map")
	}

	v, ok := f["name"]
	if !ok {
		return nil, fmt.Errorf("malformed function: no \"name\"")
	}
	c.Name, ok = v.(string)
	if !ok {
		return nil, fmt.Errorf("malformed function: \"name\" not a string")
	}

	v, ok = f["path"]
	if !ok {
		return nil, fmt.Errorf("malformed function: no \"path\"")
	}
	c.Path, ok = v.(string)
	if !ok {
		return nil, fmt.Errorf("malformed function: \"path\" not a string")
	}

	v, ok = f["arguments"]
	if !ok {
		return nil, fmt.Errorf("malformed function: no \"arguments\"")
	}
	originalArguments, ok := v.(ard.List)
	if !ok {
		return nil, fmt.Errorf("malformed function: \"arguments\" not a list")
	}

	c.Arguments = make([]Coercible, len(originalArguments))
	for index, argument := range originalArguments {
		var err error
		c.Arguments[index], err = self.NewCoercible(argument, site, source, target)
		if err != nil {
			return nil, err
		}
	}

	var err error
	c.Constraints, err = self.NewConstraintsForValue(map_, site, source, target)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (self *Function) Signature(arguments []interface{}) string {
	s := make([]string, len(arguments))
	for index, argument := range arguments {
		s[index], _ = format.EncodeJson(argument, "")
	}
	return fmt.Sprintf("%s(%s)", self.Name, strings.Join(s, ","))
}

// Coercible interface
func (self *Function) Coerce() (interface{}, error) {
	arguments, err := self.CoerceArguments()
	if err != nil {
		return nil, err
	}

	log.Infof("{evaluate} %s %s", self.Path, self.Signature(arguments))

	r, err := self.Context.CallFunction(self.Site, self.Source, self.Target, self.Name, "evaluate", arguments)
	if err != nil {
		return nil, self.WrapError(arguments, err)
	}

	// TODO: Coerce result?

	return self.Constraints.Apply(r)
}

// Coercible interface
func (self *Function) Unwrap() interface{} {
	return self.Notation
}

func (self *Function) CoerceArguments() ([]interface{}, error) {
	arguments := make([]interface{}, len(self.Arguments))
	for index, argument := range self.Arguments {
		var err error
		arguments[index], err = argument.Coerce()
		if err != nil {
			return nil, err
		}
	}
	return arguments, nil
}

func (self *Function) Validate(value interface{}, errorWhenInvalid bool) (bool, error) {
	arguments, err := self.CoerceArguments()
	if err != nil {
		return false, err
	}

	// Prepend value
	arguments = append([]interface{}{value}, arguments...)

	log.Infof("{validate} %s %s", self.Path, self.Signature(arguments))

	r, err := self.Context.CallFunction(self.Site, self.Source, self.Target, self.Name, "validate", arguments)
	if err != nil {
		return false, self.WrapError(arguments, err)
	}

	valid, ok := r.(bool)
	if !ok {
		return false, self.NewError(arguments, "\"validate\" did not return a bool")
	}

	if !valid {
		if errorWhenInvalid {
			return false, self.NewError(arguments, "")
		} else {
			return false, nil
		}
	}

	return true, nil
}
