package main

import (
	"C"
	"bytes"

	"github.com/tliron/kutil/format"
	"github.com/tliron/kutil/problems"
	"github.com/tliron/kutil/terminal"
	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/kutil/util"
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/tosca/compiler"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/puccini/tosca/parser"
)

//export Compile
func Compile(url *C.char) *C.char {
	util.ConfigureLogging(0, nil)

	buffer := bytes.NewBuffer(nil)
	terminal.Stdout = buffer

	var inputs map[string]interface{}

	var url_ urlpkg.URL
	var serviceTemplate *normal.ServiceTemplate
	var clout *cloutpkg.Clout
	var problems *problems.Problems
	var err error

	urlContext := urlpkg.NewContext()
	defer urlContext.Release()

	if url_, err = urlpkg.NewValidURL(C.GoString(url), nil, urlContext); err != nil {
		//t.Errorf("%s\n%s", err.Error(), p)
		return nil
	}

	if serviceTemplate, problems, err = parser.Parse(url_, nil, inputs); err != nil {
		//t.Errorf("%s\n%s", err.Error(), p)
		return nil
	}

	if clout, err = compiler.Compile(serviceTemplate, true); err != nil {
		//t.Errorf("%s\n%s", err.Error(), p)
		return nil
	}

	compiler.Resolve(clout, problems, urlContext, true, "yaml", true, false, false)
	if !problems.Empty() {
		//t.Errorf("%s", p)
		return nil
	}

	compiler.Coerce(clout, problems, urlContext, true, "yaml", true, false, false)
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
