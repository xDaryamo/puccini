package tosca

import (
	"strconv"
	"strings"
	"sync"

	urlpkg "github.com/tliron/kutil/url"
	"github.com/tliron/puccini/clout/js"
)

func (self *Context) ImportScriptlet(name string, path string) {
	var nativeArgumentIndexes []uint
	name, nativeArgumentIndexes = parseScriptletName(name)
	self.ScriptletNamespace.Set(name, &Scriptlet{
		Origin:                self.URL.Origin(),
		Path:                  path,
		NativeArgumentIndexes: nativeArgumentIndexes,
	})
}

func (self *Context) EmbedScriptlet(name string, scriptlet string) {
	var nativeArgumentIndexes []uint
	name, nativeArgumentIndexes = parseScriptletName(name)
	self.ScriptletNamespace.RegisterScriptlet(name, scriptlet, nativeArgumentIndexes)
}

func parseScriptletName(name string) (string, []uint) {
	// Parse optional native argument indexes specified in name
	// Notation example: my_constraint(0,1)
	var nativeArgumentIndexes []uint
	if parenthesis := strings.Index(name, "("); parenthesis != -1 {
		// We actually just assume an open parenthesis
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
		var urlContext *urlpkg.Context
		if self.Origin != nil {
			origins = []urlpkg.URL{self.Origin}
			urlContext = self.Origin.Context()
		} else {
			urlContext = urlpkg.NewContext()
			defer urlContext.Release()
		}

		url, err := urlpkg.NewValidURL(self.Path, origins, urlContext)
		if err != nil {
			return "", err
		}

		scriptlet, err := urlpkg.ReadString(url)
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
	lock      sync.RWMutex
}

func NewScriptletNamespace() *ScriptletNamespace {
	return &ScriptletNamespace{
		namespace: make(map[string]*Scriptlet),
	}
}

func (self *ScriptletNamespace) Range(f func(string, *Scriptlet) bool) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	for name, scriptlet := range self.namespace {
		if !f(name, scriptlet) {
			return
		}
	}
}

func (self *ScriptletNamespace) Lookup(name string) (*Scriptlet, bool) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	scriptlet, ok := self.namespace[name]
	return scriptlet, ok
}

func (self *ScriptletNamespace) Set(name string, scriptlet *Scriptlet) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.namespace[name] = scriptlet
}

func (self *ScriptletNamespace) RegisterScriptlet(name string, scriptlet string, nativeArgumentIndexes []uint) {
	self.Set(name, &Scriptlet{
		Scriptlet:             js.CleanupScriptlet(scriptlet),
		NativeArgumentIndexes: nativeArgumentIndexes,
	})
}

func (self *ScriptletNamespace) RegisterScriptlets(scriptlets map[string]string, nativeArgumentIndexes map[string][]uint, ignore ...string) {
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

	self.lock.Lock()
	defer self.lock.Unlock()

	namespace.lock.RLock()
	defer namespace.lock.RUnlock()

	for name, scriptlet := range namespace.namespace {
		self.namespace[name] = scriptlet
	}
}
