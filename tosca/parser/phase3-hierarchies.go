package parser

import (
	"github.com/tliron/kutil/reflection"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/puccini/tosca"
)

func (self *ServiceContext) AddHierarchies() {
	self.Context.lock.Lock()
	defer self.Context.lock.Unlock()

	self.Root.mergeHierarchies(make(tosca.HierarchyContext), self.Context.addHierarchyWork)
}

func (self *File) mergeHierarchies(hierarchyContext tosca.HierarchyContext, work reflection.EntityWork) {
	context := self.GetContext()

	self.importsLock.RLock()
	defer self.importsLock.RUnlock()

	for _, import_ := range self.Imports {
		import_.mergeHierarchies(hierarchyContext, work)
		context.Hierarchy.Merge(import_.GetContext().Hierarchy, hierarchyContext)
	}

	logHierarchies.Debugf("create: %s", context.URL.String())
	hierarchy := tosca.NewHierarchyFor(self.EntityPtr, work, hierarchyContext)
	context.Hierarchy.Merge(hierarchy, hierarchyContext)
	// TODO: do we need this?
	//context.Hierarchy.AddTo(self.EntityPtr)
}

// Print

func (self *ServiceContext) PrintHierarchies(indent int) {
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
