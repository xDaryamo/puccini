package parser

import (
	"fmt"
	"sync"

	"github.com/tliron/kutil/problems"
	"github.com/tliron/kutil/terminal"
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

func NewContext(quirks tosca.Quirks) *Context {
	return &Context{
		Quirks:          quirks,
		NamespacesWork:  NewContextualWork(logNamespaces),
		HierarchiesWork: NewContextualWork(logHierarchies),
	}
}

func (self *Context) GetProblems() *problems.Problems {
	return self.Root.GetContext().Problems
}

func (self *Context) MergeProblems() {
	// Note: This could happen many times, but because problems are de-duped, everything is OK :)
	for _, unit := range self.Units {
		self.GetProblems().Merge(unit.GetContext().Problems)
	}
}

func (self *Context) AddUnit(unit *Unit) {
	self.unitsLock.Lock()
	self.Units = append(self.Units, unit)
	self.unitsLock.Unlock()
}

func (self *Context) AddImportUnit(entityPtr tosca.EntityPtr, container *Unit, nameTransformer tosca.NameTransformer) *Unit {
	unit := NewUnit(entityPtr, container, nameTransformer)

	if container != nil {
		containerContext := container.GetContext()
		if !containerContext.HasQuirk(tosca.QuirkImportsPermissive) {
			unitContext := unit.GetContext()
			if !grammars.CompatibleGrammars(containerContext, unitContext) {
				containerContext.ReportImportIncompatible(unitContext.URL)
				return unit
			}
		}
	}

	self.AddUnit(unit)

	self.goReadImports(unit)

	return unit
}

// Print

func (self *Context) PrintImports(indent int) {
	terminal.PrintIndent(indent)
	fmt.Fprintf(terminal.Stdout, "%s\n", terminal.StyleValue(self.Root.GetContext().URL.String()))
	self.Root.PrintImports(indent, terminal.TreePrefix{})
}
