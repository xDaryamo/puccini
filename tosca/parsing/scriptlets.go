package parsing

import (
	contextpkg "context"
	"strconv"
	"strings"

	"github.com/tliron/exturl"
	"github.com/tliron/puccini/clout/js"
)

func (self *Context) ImportScriptlet(name string, path string) {
	var nativeArgumentIndexes []int
	name, nativeArgumentIndexes = parseScriptletName(name)
	self.ScriptletNamespace.Set(name, &Scriptlet{
		Origin:                self.URL.Origin(),
		Path:                  path,
		NativeArgumentIndexes: nativeArgumentIndexes,
	})
}

func (self *Context) EmbedScriptlet(name string, scriptlet string) {
	var nativeArgumentIndexes []int
	name, nativeArgumentIndexes = parseScriptletName(name)
	self.ScriptletNamespace.RegisterScriptlet(name, scriptlet, nativeArgumentIndexes)
}

//
// Scriptlet
//

type Scriptlet struct {
	Origin                exturl.URL `json:"origin" yaml:"origin"`
	Path                  string     `json:"path" yaml:"path"`
	Scriptlet             string     `json:"scriptlet" yaml:"scriptlet"`
	NativeArgumentIndexes []int      `json:"nativeArgumentIndexes" yaml:"nativeArgumentIndexes"`
}

func (self *Scriptlet) Read(context contextpkg.Context) (string, error) {
	if self.Path != "" {
		var origins []exturl.URL
		var urlContext *exturl.Context
		if self.Origin != nil {
			origins = []exturl.URL{self.Origin}
			urlContext = self.Origin.Context()
		} else {
			urlContext = exturl.NewContext()
			defer urlContext.Release()
		}

		url, err := urlContext.NewValidURL(context, self.Path, origins)
		if err != nil {
			return "", err
		}

		scriptlet, err := exturl.ReadString(context, url)
		if err != nil {
			return "", err
		}

		return scriptlet, nil
	}

	return self.Scriptlet, nil
}

//
// ScriptletNamespace
//

type ScriptletNamespace struct {
	namespace map[string]*Scriptlet
}

func NewScriptletNamespace() *ScriptletNamespace {
	return &ScriptletNamespace{
		namespace: make(map[string]*Scriptlet),
	}
}

func (self *ScriptletNamespace) Range(f func(string, *Scriptlet) bool) {
	namespace := make(map[string]*Scriptlet)
	for name, scriptlet := range self.namespace {
		namespace[name] = scriptlet
	}

	for name, scriptlet := range namespace {
		if !f(name, scriptlet) {
			return
		}
	}
}

func (self *ScriptletNamespace) Lookup(name string) (*Scriptlet, bool) {
	scriptlet, ok := self.namespace[name]
	return scriptlet, ok
}

func (self *ScriptletNamespace) Set(name string, scriptlet *Scriptlet) {
	self.namespace[name] = scriptlet
}

func (self *ScriptletNamespace) RegisterScriptlet(name string, scriptlet string, nativeArgumentIndexes []int) {
	self.Set(name, &Scriptlet{
		Scriptlet:             js.CleanupScriptlet(scriptlet),
		NativeArgumentIndexes: nativeArgumentIndexes,
	})
}

func (self *ScriptletNamespace) RegisterScriptlets(scriptlets map[string]string, nativeArgumentIndexes map[string][]int, ignore ...string) {
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

func (self *ScriptletNamespace) Merge(namespace *ScriptletNamespace) {
	if self == namespace {
		return
	}

	for name, scriptlet := range namespace.namespace {
		self.namespace[name] = scriptlet
	}
}

// Utils

func parseScriptletName(name string) (string, []int) {
	// Parse optional native argument indexes specified in name
	// Notation example: my_constraint(0,1)
	var nativeArgumentIndexes []int
	if parenthesis := strings.Index(name, "("); parenthesis != -1 {
		// We actually just assume an open parenthesis
		split := strings.Split(name[parenthesis+1:len(name)-1], ",")
		name = name[:parenthesis]
		for _, s := range split {
			if index, err := strconv.ParseInt(s, 10, 32); err != nil {
				nativeArgumentIndexes = append(nativeArgumentIndexes, int(index))
			}
		}
	}
	return name, nativeArgumentIndexes
}
