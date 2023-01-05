package main

import (
	"C"
	"bytes"
	"errors"
	"strings"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/problems"
	"github.com/tliron/kutil/transcribe"
	urlpkg "github.com/tliron/kutil/url"
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/clout/js"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/parser"
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

	var quirks_ tosca.Quirks

	if data, err := yamlkeys.DecodeAll(strings.NewReader(C.GoString(quirks))); err == nil {
		for _, data_ := range data {
			if list, ok := data_.(ard.List); ok {
				for _, value := range list {
					if value_, ok := value.(string); ok {
						quirks_ = append(quirks_, tosca.Quirk(value_))
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

	var url_ urlpkg.URL
	var serviceTemplate *normal.ServiceTemplate
	var clout *cloutpkg.Clout
	var problems *problems.Problems
	var err error

	urlContext := urlpkg.NewContext()
	defer urlContext.Release()

	if url_, err = urlpkg.NewValidURL(C.GoString(url), nil, urlContext); err != nil {
		return result(nil, nil, err)
	}

	if _, serviceTemplate, problems, err = parserContext.Parse(url_, nil, nil, quirks_, inputs_); err != nil {
		return result(nil, problems, err)
	}

	if clout, err = serviceTemplate.Compile(); err != nil {
		return result(clout, problems, err)
	}

	if resolve != 0 {
		js.Resolve(clout, problems, urlContext, true, "yaml", true, false)
		if !problems.Empty() {
			return result(clout, problems, nil)
		}
	}

	if coerce != 0 {
		js.Coerce(clout, problems, urlContext, true, "yaml", true, false)
		if !problems.Empty() {
			return result(clout, problems, nil)
		}
	}

	return result(clout, problems, nil)
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
