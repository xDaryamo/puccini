package tosca

import (
	"strconv"
	"strings"

	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/url"
)

func (self *Context) ImportScript(name string, path string) {
	var nativeArgumentIndexes []uint
	name, nativeArgumentIndexes = parseScriptName(name)
	self.ScriptNamespace[name] = &Script{
		Origin:                self.URL.Origin(),
		Path:                  path,
		NativeArgumentIndexes: nativeArgumentIndexes,
	}
}

func (self *Context) SourceScript(name string, sourceCode string) {
	var nativeArgumentIndexes []uint
	name, nativeArgumentIndexes = parseScriptName(name)
	self.ScriptNamespace[name] = &Script{
		SourceCode:            js.Cleanup(sourceCode),
		NativeArgumentIndexes: nativeArgumentIndexes,
	}
}

func parseScriptName(name string) (string, []uint) {
	// Parse optional native argument indexes
	// e.g.: my_constraint(0,1)
	var nativeArgumentIndexes []uint
	if parentheses := strings.Index(name, "("); parentheses != -1 {
		split := strings.Split(name[parentheses+1:len(name)-1], ",")
		name = name[:parentheses]
		for _, s := range split {
			if index, err := strconv.ParseUint(s, 10, 32); err != nil {
				nativeArgumentIndexes = append(nativeArgumentIndexes, uint(index))
			}
		}
	}
	return name, nativeArgumentIndexes
}

//
// Script
//

type Script struct {
	Origin                url.URL `json:"origin" yaml:"origin"`
	Path                  string  `json:"path" yaml:"path"`
	SourceCode            string  `json:"sourceCode" yaml:"sourceCode"`
	NativeArgumentIndexes []uint  `json:"nativeArgumentIndexes" yaml:"nativeArgumentIndexes"`
}

func (self *Script) GetSourceCode() (string, error) {
	if self.Path != "" {
		var origins []url.URL
		if self.Origin != nil {
			origins = []url.URL{self.Origin}
		}

		url_, err := url.NewValidURL(self.Path, origins)
		if err != nil {
			return "", err
		}

		sourceCode, err := url.Read(url_)
		if err != nil {
			return "", err
		}

		return js.Cleanup(sourceCode), nil
	}

	return self.SourceCode, nil
}

//
// ScriptNamespace
//

type ScriptNamespace map[string]*Script

func (self ScriptNamespace) Merge(namespace ScriptNamespace) {
	for name, script := range namespace {
		self[name] = script
	}
}
