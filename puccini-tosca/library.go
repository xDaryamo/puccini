package main

import (
	"C"
	"bytes"
	"errors"
	"strings"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/format"
	"github.com/tliron/kutil/logging"
	"github.com/tliron/kutil/problems"
	urlpkg "github.com/tliron/kutil/url"
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/tosca/compiler"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/parser"
	"github.com/tliron/yamlkeys"
)

//export Compile
func Compile(url *C.char, inputs *C.char) *C.char {
	logging.Configure(0, nil)

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

	if _, serviceTemplate, problems, err = parser.Parse(url_, nil, inputs_); err != nil {
		return result(nil, problems, err)
	}

	if clout, err = compiler.Compile(serviceTemplate, true); err != nil {
		return result(clout, problems, err)
	}

	compiler.Resolve(clout, problems, urlContext, true, "yaml", true, false, false)
	if !problems.Empty() {
		return result(clout, problems, nil)
	}

	compiler.Coerce(clout, problems, urlContext, true, "yaml", true, false, false)
	if !problems.Empty() {
		return result(clout, problems, nil)
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
	format.WriteYAML(result, buffer, "  ", true)
	return C.CString(buffer.String())
}
