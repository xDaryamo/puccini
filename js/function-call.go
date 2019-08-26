package js

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/format"
)

//
// FunctionCall
//

type FunctionCall struct {
	Context     *CloutContext `json:"-" yaml:"-"`
	Name        string        `json:"functionCall" yaml:"functionCall"`
	URL         string        `json:"url" yaml:"url"`
	Path        string        `json:"path" yaml:"path"`
	Location    string        `json:"location" yaml:"location"`
	Arguments   []Coercible   `json:"arguments" yaml:"arguments"`
	Constraints Constraints   `json:"constraints" yaml:"constraints"`
	Site        interface{}   `json:"-" yaml:"-"`
	Source      interface{}   `json:"-" yaml:"-"`
	Target      interface{}   `json:"-" yaml:"-"`
	Notation    ard.Map       `json:"-" yaml:"-"`
}

func (self *CloutContext) NewFunctionCall(data interface{}, site interface{}, source interface{}, target interface{}) (*FunctionCall, error) {
	map_, ok := data.(ard.Map)
	if !ok {
		return nil, errors.New("not a function call")
	}

	functionCall, ok := map_["functionCall"]
	if !ok {
		return nil, errors.New("not a function call")
	}

	c := FunctionCall{
		Context:  self,
		Site:     site,
		Source:   source,
		Target:   target,
		Notation: map_,
	}

	f, ok := functionCall.(ard.Map)
	if !ok {
		return nil, errors.New("malformed function call: not a map")
	}

	v, ok := f["name"]
	if !ok {
		return nil, errors.New("malformed function call: no \"name\"")
	}
	c.Name, ok = v.(string)
	if !ok {
		return nil, errors.New("malformed function call: \"name\" not a string")
	}

	if v, ok = f["url"]; ok {
		c.URL, ok = v.(string)
		if !ok {
			return nil, errors.New("malformed function call: \"url\" not a string")
		}
	}

	if v, ok = f["path"]; ok {
		c.Path, ok = v.(string)
		if !ok {
			return nil, errors.New("malformed function call: \"path\" not a string")
		}
	}

	if v, ok = f["location"]; ok {
		c.Location, ok = v.(string)
		if !ok {
			return nil, errors.New("malformed function call: \"location\" not a string")
		}
	}

	v, ok = f["arguments"]
	if !ok {
		return nil, errors.New("malformed function call: no \"arguments\"")
	}
	originalArguments, ok := v.(ard.List)
	if !ok {
		return nil, errors.New("malformed function call: \"arguments\" not a list")
	}

	c.Arguments = make([]Coercible, len(originalArguments))
	for index, argument := range originalArguments {
		var err error
		if c.Arguments[index], err = self.NewCoercible(argument, site, source, target); err != nil {
			return nil, err
		}
	}

	var err error
	if c.Constraints, err = self.NewConstraintsForValue(map_, site, source, target); err != nil {
		return nil, err
	}

	return &c, nil
}

func (self *FunctionCall) Signature(arguments []interface{}) string {
	s := make([]string, len(arguments))
	for index, argument := range arguments {
		s[index], _ = format.EncodeJson(argument, "")
	}
	return fmt.Sprintf("%s(%s)", self.Name, strings.Join(s, ","))
}

// Coercible interface
func (self *FunctionCall) Coerce() (interface{}, error) {
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
func (self *FunctionCall) Unwrap() interface{} {
	return self.Notation
}

func (self *FunctionCall) CoerceArguments() ([]interface{}, error) {
	arguments := make([]interface{}, len(self.Arguments))
	for index, argument := range self.Arguments {
		var err error
		if arguments[index], err = argument.Coerce(); err != nil {
			return nil, err
		}
	}
	return arguments, nil
}

func (self *FunctionCall) Validate(value interface{}, errorWhenInvalid bool) (bool, error) {
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

	if valid, ok := r.(bool); ok {
		if !valid {
			if errorWhenInvalid {
				return false, self.NewError(arguments, "", nil)
			} else {
				return false, nil
			}
		}
	} else {
		return false, self.NewError(arguments, "\"validate\" did not return a bool", nil)
	}

	return true, nil
}
