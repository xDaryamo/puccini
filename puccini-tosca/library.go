package main

import (
	"C"
	"bytes"

	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/tosca/compiler"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/parser"
	"github.com/tliron/puccini/tosca/problems"
)

//export Compile
func Compile(url *C.char) *C.char {
	common.ConfigureLogging(0, nil)

	buffer := bytes.NewBuffer(nil)
	format.Stdout = buffer

	var inputs map[string]interface{}

	var serviceTemplate *normal.ServiceTemplate
	var clout_ *clout.Clout
	var problems *problems.Problems
	var err error

	if serviceTemplate, problems, err = parser.Parse(C.GoString(url), nil, inputs); err != nil {
		//t.Errorf("%s\n%s", err.Error(), p)
		return nil
	}

	if clout_, err = compiler.Compile(serviceTemplate); err != nil {
		//t.Errorf("%s\n%s", err.Error(), p)
		return nil
	}

	compiler.Resolve(clout_, problems, "yaml", true)
	if !problems.Empty() {
		//t.Errorf("%s", p)
		return nil
	}

	compiler.Coerce(clout_, problems, "yaml", true)
	if !problems.Empty() {
		//t.Errorf("%s", p)
		return nil
	}

	format.WriteYaml(clout_, buffer, "  ")

	return C.CString(buffer.String())
}
