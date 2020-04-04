package parser

import (
	"fmt"
	"sync"

	"github.com/tliron/puccini/common/problems"
	"github.com/tliron/puccini/common/terminal"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars"
)

type Context struct {
	Root            *Unit
	Quirks          tosca.Quirks
	Units           Units
	Parsing         sync.Map
	WaitGroup       sync.WaitGroup
	NamespacesWork  *ContextualWork
	HierarchiesWork *ContextualWork

	unitsLock sync.Mutex
}

func NewContext(quirks tosca.Quirks) Context {
	return Context{
		Quirks:          quirks,
		NamespacesWork:  NewContextualWork("namespaces"),
		HierarchiesWork: NewContextualWork("hierarchies"),
	}
}

func (self *Context) GetProblems() *problems.Problems {
	return self.Root.GetContext().Problems
}

func (self *Context) AddUnit(entityPtr tosca.EntityPtr, container *Unit, nameTransformer tosca.NameTransformer) *Unit {
	unit := NewUnit(entityPtr, container, nameTransformer)

	if container != nil {
		containerContext := container.GetContext()
		if !containerContext.HasQuirk(tosca.QuirkImportsImperssive) {
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

	self.unitsLock.Lock()
	self.Units = append(self.Units, unit)
	self.unitsLock.Unlock()

	self.goReadImports(unit)

	return unit
}

// Print

func (self *Context) PrintImports(indent int) {
	terminal.PrintIndent(indent)
	fmt.Fprintf(terminal.Stdout, "%s\n", terminal.ColorValue(self.Root.GetContext().URL.String()))
	self.Root.PrintImports(indent, terminal.TreePrefix{})
}
