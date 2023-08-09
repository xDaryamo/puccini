package parser

import (
	"github.com/tliron/kutil/reflection"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/puccini/tosca/parsing"
)

func (self *Context) AddHierarchies() {
	self.Parser.lock.Lock()
	defer self.Parser.lock.Unlock()

	self.Root.mergeHierarchies(make(parsing.HierarchyContext), self.Parser.addHierarchyWork)
}

func (self *File) mergeHierarchies(hierarchyContext parsing.HierarchyContext, work reflection.EntityWork) {
	context := self.GetContext()

	self.importsLock.RLock()
	defer self.importsLock.RUnlock()

	for _, import_ := range self.Imports {
		import_.mergeHierarchies(hierarchyContext, work)
		context.Hierarchy.Merge(import_.GetContext().Hierarchy, hierarchyContext)
	}

	logHierarchies.Debugf("create: %s", context.URL.String())
	hierarchy := parsing.NewHierarchyFor(self.EntityPtr, work, hierarchyContext)
	context.Hierarchy.Merge(hierarchy, hierarchyContext)
	// TODO: do we need this?
	//context.Hierarchy.AddTo(self.EntityPtr)
}

// Print

func (self *Context) PrintHierarchies(indent int) {
	self.filesLock.RLock()
	defer self.filesLock.RUnlock()

	for _, file := range self.Files {
		context := file.GetContext()
		if !context.Hierarchy.Empty() {
			terminal.PrintIndent(indent)
			terminal.Printf("%s\n", self.Stylist.Value(context.URL.String()))
			context.Hierarchy.Print(indent)
		}
	}
}
