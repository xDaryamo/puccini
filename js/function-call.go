package js

import (
	"fmt"
	"strings"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/format"
)

//
// FunctionCallContext
//

type FunctionCallContext struct {
	Site   interface{}
	Source interface{}
	Target interface{}
}

func (self FunctionCallContext) Map() map[string]interface{} {
	return map[string]interface{}{
		"site":   self.Site,
		"source": self.Source,
		"target": self.Target,
	}
}

//
// FunctionCall
//

type FunctionCall struct {
	RuntimeContext      *RuntimeContext     `json:"-" yaml:"-"`
	FunctionCallContext FunctionCallContext `json:"-" yaml:"-"`

	Name        string      `json:"functionCall" yaml:"functionCall"`
	Arguments   []Coercible `json:"arguments" yaml:"arguments"`
	URL         string      `json:"url" yaml:"url"`
	Path        string      `json:"path" yaml:"path"`
	Location    string      `json:"location" yaml:"location"`
	Constraints Constraints `json:"constraints" yaml:"constraints"`

	Notation ard.StringMap `json:"-" yaml:"-"`
}

func (self *RuntimeContext) NewFunctionCall(map_ ard.StringMap, functionCallContext FunctionCallContext) (*FunctionCall, error) {
	coercible := FunctionCall{
		RuntimeContext:      self,
		FunctionCallContext: functionCallContext,
	}

	if data, ok := map_["name"]; ok {
		if coercible.Name, ok = data.(string); !ok {
			return nil, fmt.Errorf("malformed function call, \"name\" not a string: %T", data)
		}
	} else {
		return nil, fmt.Errorf("malformed function call, no \"name\": %v", map_)
	}

	if data, ok := map_["arguments"]; ok {
		if originalArguments, ok := data.(ard.List); ok {
			coercible.Arguments = make([]Coercible, len(originalArguments))
			for index, argument := range originalArguments {
				var err error
				if coercible.Arguments[index], err = self.NewCoercible(argument, functionCallContext); err != nil {
					return nil, err
				}
			}
		} else {
			return nil, fmt.Errorf("malformed function call, \"arguments\" not a list: %T", data)
		}
	} else {
		return nil, fmt.Errorf("malformed function call, no \"arguments\": %v", map_)
	}

	if data, ok := map_["url"]; ok {
		if coercible.URL, ok = data.(string); !ok {
			return nil, fmt.Errorf("malformed function call, \"url\" not a string: %T", data)
		}
	}

	if data, ok := map_["path"]; ok {
		if coercible.Path, ok = data.(string); !ok {
			return nil, fmt.Errorf("malformed function call, \"path\" not a string: %T", data)
		}
	}

	if data, ok := map_["location"]; ok {
		if coercible.Location, ok = data.(string); !ok {
			return nil, fmt.Errorf("malformed function call, \"location\" not a string: %T", data)
		}
	}

	return &coercible, nil
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

	r, err := self.RuntimeContext.CallFunction(self.Name, "evaluate", arguments, self.FunctionCallContext)
	if err != nil {
		return nil, self.WrapError(arguments, err)
	}

	// TODO: Coerce result?

	return self.Constraints.Apply(r)
}

// Coercible interface
func (self *FunctionCall) SetConstraints(constraints Constraints) {
	self.Constraints = constraints
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

	// Prepend value to be first argument
	arguments = append([]interface{}{value}, arguments...)

	log.Infof("{validate} %s %s", self.Path, self.Signature(arguments))

	r, err := self.RuntimeContext.CallFunction(self.Name, "validate", arguments, self.FunctionCallContext)
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
