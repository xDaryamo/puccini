package main

import (
	"C"
	"bytes"
	contextpkg "context"
	"errors"
	"strings"

	"github.com/tliron/exturl"
	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/problems"
	"github.com/tliron/kutil/transcribe"
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/clout/js"
	"github.com/tliron/puccini/tosca/parser"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

var parserContext = parser.NewContext()

//export Compile
func Compile(url *C.char, inputs *C.char, quirks *C.char, resolve C.char, coerce C.char) *C.char {
	inputs_ := make(map[string]ard.Value)

	if data, err := yamlkeys.DecodeAll(strings.NewReader(C.GoString(inputs))); err == nil {
		for _, data_ := range data {
			if map_, ok := data_.(ard.Map); ok {
				for key, value := range map_ {
					inputs_[yamlkeys.KeyString(key)] = value
				}
			} else {
				return result(nil, nil, errors.New("malformed inputs"))
			}
		}
	} else {
		return result(nil, nil, err)
	}

	var quirks_ parsing.Quirks

	if data, err := yamlkeys.DecodeAll(strings.NewReader(C.GoString(quirks))); err == nil {
		for _, data_ := range data {
			if list, ok := data_.(ard.List); ok {
				for _, value := range list {
					if value_, ok := value.(string); ok {
						quirks_ = append(quirks_, parsing.Quirk(value_))
					} else {
						return result(nil, nil, errors.New("malformed quirk"))
					}
				}
			} else {
				return result(nil, nil, errors.New("malformed quirks"))
			}
		}
	} else {
		return result(nil, nil, err)
	}

	var url_ exturl.URL
	var result_ parser.Result
	var clout *cloutpkg.Clout
	var err error

	urlContext := exturl.NewContext()
	defer urlContext.Release()
	context := contextpkg.TODO()

	if url_, err = urlContext.NewValidURL(context, C.GoString(url), nil); err != nil {
		return result(nil, nil, err)
	}

	if result_, err = parserContext.Parse(context, parser.ParseContext{URL: url_, Quirks: quirks_, Inputs: inputs_}); err != nil {
		return result(nil, result_.Problems, err)
	}

	if clout, err = result_.NormalServiceTemplate.Compile(); err != nil {
		return result(clout, result_.Problems, err)
	}

	execContext := js.ExecContext{
		Clout:      clout,
		Problems:   result_.Problems,
		URLContext: urlContext,
		History:    true,
		Format:     "yaml",
		Strict:     true,
		Pretty:     false,
	}

	if resolve != 0 {
		execContext.Resolve()
		if !result_.Problems.Empty() {
			return result(clout, result_.Problems, nil)
		}
	}

	if coerce != 0 {
		execContext.Coerce()
		if !result_.Problems.Empty() {
			return result(clout, result_.Problems, nil)
		}
	}

	return result(clout, result_.Problems, nil)
}

func result(clout *cloutpkg.Clout, problems *problems.Problems, err error) *C.char {
	result := make(ard.StringMap)
	if clout != nil {
		result["clout"] = clout
	}
	if (problems != nil) && !problems.Empty() {
		result["problems"] = problems.Problems
	}
	if err != nil {
		result["error"] = err.Error()
	}

	buffer := bytes.NewBuffer(nil)
	transcribe.WriteYAML(result, buffer, "  ", true) // TODO: err
	return C.CString(buffer.String())
}
