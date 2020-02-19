package tosca

import (
	"strconv"
	"strings"

	"github.com/tliron/puccini/clout/js"
	urlpkg "github.com/tliron/puccini/url"
)

func (self *Context) ImportScriptlet(name string, path string) {
	var nativeArgumentIndexes []uint
	name, nativeArgumentIndexes = parseScriptletName(name)
	self.ScriptletNamespace[name] = &Scriptlet{
		Origin:                self.URL.Origin(),
		Path:                  path,
		NativeArgumentIndexes: nativeArgumentIndexes,
	}
}

func (self *Context) EmbedScriptlet(name string, scriptlet string) {
	var nativeArgumentIndexes []uint
	name, nativeArgumentIndexes = parseScriptletName(name)
	self.ScriptletNamespace[name] = &Scriptlet{
		Scriptlet:             js.CleanupScriptlet(scriptlet),
		NativeArgumentIndexes: nativeArgumentIndexes,
	}
}

func parseScriptletName(name string) (string, []uint) {
	// Parse optional native argument indexes specified in name
	// Notation example: my_constraint(0,1)
	var nativeArgumentIndexes []uint
	if parenthesis := strings.Index(name, "("); parenthesis != -1 {
		// We actually just assume an open paranthesis
		split := strings.Split(name[parenthesis+1:len(name)-1], ",")
		name = name[:parenthesis]
		for _, s := range split {
			if index, err := strconv.ParseUint(s, 10, 32); err != nil {
				nativeArgumentIndexes = append(nativeArgumentIndexes, uint(index))
			}
		}
	}
	return name, nativeArgumentIndexes
}

//
// Scriptlet
//

type Scriptlet struct {
	Origin                urlpkg.URL `json:"origin" yaml:"origin"`
	Path                  string     `json:"path" yaml:"path"`
	Scriptlet             string     `json:"scriptlet" yaml:"scriptlet"`
	NativeArgumentIndexes []uint     `json:"nativeArgumentIndexes" yaml:"nativeArgumentIndexes"`
}

func (self *Scriptlet) Read() (string, error) {
	if self.Path != "" {
		var origins []urlpkg.URL
		if self.Origin != nil {
			origins = []urlpkg.URL{self.Origin}
		}

		url, err := urlpkg.NewValidURL(self.Path, origins)
		if err != nil {
			return "", err
		}

		scriptlet, err := urlpkg.Read(url)
		if err != nil {
			return "", err
		}

		return js.CleanupScriptlet(scriptlet), nil
	}

	return self.Scriptlet, nil
}

//
// ScriptletNamespace
//

type ScriptletNamespace map[string]*Scriptlet

func (self ScriptletNamespace) RegisterScriptlet(name string, scriptlet string, nativeArgumentIndexes []uint) {
	self[name] = &Scriptlet{
		Scriptlet:             js.CleanupScriptlet(scriptlet),
		NativeArgumentIndexes: nativeArgumentIndexes,
	}
}

func (self ScriptletNamespace) RegisterScriptlets(scriptlets map[string]string, nativeArgumentIndexes map[string][]uint, ignore ...string) {
	for name, scriptlet := range scriptlets {
		var ignore_ bool
		for _, ignore__ := range ignore {
			if name == ignore__ {
				ignore_ = true
				break
			}
		}

		if ignore_ {
			continue
		}

		self.RegisterScriptlet(name, scriptlet, nativeArgumentIndexes[name])
	}
}

func (self ScriptletNamespace) Merge(from ScriptletNamespace) {
	for name, scriptlet := range from {
		self[name] = scriptlet
	}
}
