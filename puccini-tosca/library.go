package main

import (
	"C"
	"bytes"

	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/common/format"
	"github.com/tliron/puccini/common/terminal"
	"github.com/tliron/puccini/tosca/compiler"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/parser"
	"github.com/tliron/puccini/tosca/problems"
)

//export Compile
func Compile(url *C.char) *C.char {
	common.ConfigureLogging(0, nil)

	buffer := bytes.NewBuffer(nil)
	terminal.Stdout = buffer

	var inputs map[string]interface{}

	var serviceTemplate *normal.ServiceTemplate
	var clout *cloutpkg.Clout
	var problems *problems.Problems
	var err error

	if serviceTemplate, problems, err = parser.Parse(C.GoString(url), nil, inputs); err != nil {
		//t.Errorf("%s\n%s", err.Error(), p)
		return nil
	}

	if clout, err = compiler.Compile(serviceTemplate); err != nil {
		//t.Errorf("%s\n%s", err.Error(), p)
		return nil
	}

	compiler.Resolve(clout, problems, "yaml", true, false)
	if !problems.Empty() {
		//t.Errorf("%s", p)
		return nil
	}

	compiler.Coerce(clout, problems, "yaml", true, false)
	if !problems.Empty() {
		//t.Errorf("%s", p)
		return nil
	}

	ard, err := clout.ARD()
	if err != nil {
		return nil
	}

	format.WriteYAML(ard, buffer, "  ", true)

	return C.CString(buffer.String())
}
