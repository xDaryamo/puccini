package parser

import (
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/puccini/tosca"
)

func (self *ServiceContext) AddHierarchies() {
	self.Context.entitiesLock.Lock()
	defer self.Context.entitiesLock.Unlock()

	self.Root.MergeHierarchies(make(tosca.HierarchyContext), self.Context.addHierarchyWork)
}

func (self *Unit) MergeHierarchies(hierarchyContext tosca.HierarchyContext, work tosca.EntityWork) {
	context := self.GetContext()

	self.importsLock.RLock()
	defer self.importsLock.RUnlock()

	for _, import_ := range self.Imports {
		import_.MergeHierarchies(hierarchyContext, work)
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
	self.unitsLock.RLock()
	defer self.unitsLock.RUnlock()

	for _, unit := range self.Units {
		context := unit.GetContext()
		if !context.Hierarchy.Empty() {
			terminal.PrintIndent(indent)
			terminal.Printf("%s\n", self.Stylist.Value(context.URL.String()))
			context.Hierarchy.Print(indent)
		}
	}
}
