package parser

import (
	"sync"

	"github.com/tliron/kutil/problems"
	"github.com/tliron/kutil/reflection"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars"
)

//
// Context
//

type Context struct {
	readCache          sync.Map // entityPtr or Promise
	lookupFieldsWork   reflection.EntityWork
	addHierarchyWork   reflection.EntityWork
	getInheritTaskWork reflection.EntityWork
	renderWork         reflection.EntityWork
	lock               util.RWLocker
}

func NewContext() *Context {
	return &Context{
		lookupFieldsWork:   make(reflection.EntityWork),
		addHierarchyWork:   make(reflection.EntityWork),
		getInheritTaskWork: make(reflection.EntityWork),
		renderWork:         make(reflection.EntityWork),
		lock:               util.NewDefaultRWLocker(),
	}
}

//
// ServiceContext
//

type ServiceContext struct {
	Context *Context
	Root    *Unit
	Stylist *terminal.Stylist
	Quirks  tosca.Quirks
	Units   Units

	readWork  sync.WaitGroup
	unitsLock util.RWLocker
}

func (self *Context) NewServiceContext(stylist *terminal.Stylist, quirks tosca.Quirks) *ServiceContext {
	return &ServiceContext{
		Context:   self,
		Stylist:   stylist,
		Quirks:    quirks,
		unitsLock: util.NewDebugRWLocker(),
	}
}

func (self *ServiceContext) GetProblems() *problems.Problems {
	return self.Root.GetContext().Problems
}

func (self *ServiceContext) MergeProblems() {
	self.unitsLock.RLock()
	defer self.unitsLock.RUnlock()

	// Note: This could happen many times, but because problems are de-duped, everything is OK :)
	for _, unit := range self.Units {
		self.GetProblems().Merge(unit.GetContext().Problems)
	}
}

func (self *ServiceContext) AddUnit(unit *Unit) {
	self.unitsLock.Lock()
	defer self.unitsLock.Unlock()

	self.Units = append(self.Units, unit)
}

func (self *ServiceContext) AddImportUnit(entityPtr tosca.EntityPtr, container *Unit, nameTransformer tosca.NameTransformer) *Unit {
	unit := NewUnit(entityPtr, container, nameTransformer)

	if container != nil {
		containerContext := container.GetContext()
		if !containerContext.HasQuirk(tosca.QuirkImportsVersionPermissive) {
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

func (self *ServiceContext) PrintImports(indent int) {
	terminal.PrintIndent(indent)
	terminal.Printf("%s\n", terminal.Stylize.Value(self.Root.GetContext().URL.String()))
	self.Root.PrintImports(indent, terminal.TreePrefix{})
}
