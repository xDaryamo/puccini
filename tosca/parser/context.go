package parser

import (
	"fmt"
	"sync"

	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/problems"
)

type Context struct {
	Root     *Unit
	Problems problems.Problems
	Quirks   []string
	Units    Units
	Parsing  sync.Map
	WG       sync.WaitGroup
	Locker   sync.Mutex
}

func NewContext(quirks []string) Context {
	return Context{
		Quirks: quirks,
	}
}

func (self *Context) AddUnit(unit *Unit) {
	self.Locker.Lock()
	self.Units = append(self.Units, unit)
	self.Locker.Unlock()
}

func (self *Context) AddUnitFor(entityPtr interface{}, container *Unit, nameTransformer tosca.NameTransformer) {
	unit := NewUnit(entityPtr, container, nameTransformer)
	if container != nil {
		// Merge problems into container
		// Note: In concurrent uses of the same unit this can lead to the same problems being merged more than once
		// But because duplicate problems are not merged, everything is OK :)
		container.GetContext().Problems.Merge(unit.GetContext().Problems)
	}
	self.goReadImports(unit)
}

// Print

func (self *Context) PrintImports(indent int) {
	format.PrintIndent(indent)
	fmt.Fprintf(format.Stdout, "%s\n", format.ColorValue(self.Root.GetContext().URL.String()))
	self.Root.PrintImports(indent, format.TreePrefix{})
}
