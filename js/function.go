package js

import (
	"fmt"
	"strings"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/format"
)

//
// Function
//

type Function struct {
	Clout       *clout.Clout `json:"-" yaml:"-"`
	Name        string       `json:"name" yaml:"name"`
	Path        string       `json:"path" yaml:"path"`
	Arguments   []Coercible  `json:"arguments" yaml:"arguments"`
	Constraints Constraints  `json:"constraints" yaml:"constraints"`
	Site        interface{}  `json:"-" yaml:"-"`
	Source      interface{}  `json:"-" yaml:"-"`
	Target      interface{}  `json:"-" yaml:"-"`
}

func NewFunction(data interface{}, site interface{}, source interface{}, target interface{}, c *clout.Clout) (*Function, error) {
	map_, ok := data.(ard.Map)
	if !ok {
		return nil, fmt.Errorf("not a function")
	}

	function, ok := map_["function"]
	if !ok {
		return nil, fmt.Errorf("not a function")
	}

	self := Function{Clout: c, Site: site, Source: source, Target: target}

	f, ok := function.(ard.Map)
	if !ok {
		return nil, fmt.Errorf("malformed function: not a map")
	}

	v, ok := f["name"]
	if !ok {
		return nil, fmt.Errorf("malformed function: no \"name\"")
	}
	self.Name, ok = v.(string)
	if !ok {
		return nil, fmt.Errorf("malformed function: \"name\" not a string")
	}

	v, ok = f["path"]
	if !ok {
		return nil, fmt.Errorf("malformed function: no \"path\"")
	}
	self.Path, ok = v.(string)
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

	self.Arguments = make([]Coercible, len(originalArguments))
	for index, argument := range originalArguments {
		var err error
		self.Arguments[index], err = NewCoercible(argument, site, source, target, c)
		if err != nil {
			return nil, err
		}
	}

	var err error
	self.Constraints, err = NewConstraints(map_, c)
	if err != nil {
		return nil, err
	}

	return &self, nil
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

	r, err := CallClout(self.Clout, self.Site, self.Source, self.Target, self.Name, "evaluate", arguments)
	if err != nil {
		return nil, self.WrapError(arguments, err)
	}

	// TODO: Coerce result?

	return self.Constraints.Apply(r)
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

func (self *Function) Validate(value interface{}) (bool, error) {
	arguments, err := self.CoerceArguments()
	if err != nil {
		return false, err
	}

	// Prepend value
	arguments = append([]interface{}{value}, arguments...)

	log.Infof("{validate} %s %s", self.Path, self.Signature(arguments))

	r, err := CallClout(self.Clout, self.Site, self.Source, self.Target, self.Name, "validate", arguments)
	if err != nil {
		return false, self.WrapError(arguments, err)
	}

	valid, ok := r.(bool)
	if !ok {
		return false, self.NewError(arguments, "\"validate\" did not return a bool")
	}

	if !valid {
		return false, self.NewError(arguments, "")
	}

	return true, nil
}
