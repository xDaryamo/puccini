package parser

import (
	"fmt"
	"sync"

	"github.com/tliron/puccini/common/terminal"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars"
	"github.com/tliron/puccini/tosca/problems"
)

type Context struct {
	Root      *Unit
	Quirks    []string
	Units     Units
	Parsing   sync.Map
	WaitGroup sync.WaitGroup
	Locker    sync.Mutex
}

func NewContext(quirks []string) Context {
	return Context{
		Quirks: quirks,
	}
}

func (self *Context) GetProblems() *problems.Problems {
	return self.Root.GetContext().Problems
}

func (self *Context) AddUnit(entityPtr interface{}, container *Unit, nameTransformer tosca.NameTransformer) *Unit {
	unit := NewUnit(entityPtr, container, nameTransformer)

	if container != nil {
		containerContext := container.GetContext()
		if !containerContext.HasQuirk("imports.permissive") {
			unitContext := unit.GetContext()
			if !grammars.CompatibleGrammars(containerContext, unitContext) {
				containerContext.ReportImportIncompatible(unitContext.URL)
				return unit
			}
		}

		// Merge problems into container
		// Note: This happens every time the same unit is imported,
		// so it could be that that problems are merged multiple times,
		// but because problems are de-duped, everything is OK :)
		container.GetContext().Problems.Merge(unit.GetContext().Problems)
	}

	self.Locker.Lock()
	self.Units = append(self.Units, unit)
	self.Locker.Unlock()

	self.goReadImports(unit)

	return unit
}

// Print

func (self *Context) PrintImports(indent int) {
	terminal.PrintIndent(indent)
	fmt.Fprintf(terminal.Stdout, "%s\n", terminal.ColorValue(self.Root.GetContext().URL.String()))
	self.Root.PrintImports(indent, terminal.TreePrefix{})
}
