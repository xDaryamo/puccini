package parser

import (
	"fmt"

	"github.com/tliron/kutil/terminal"
	"github.com/tliron/puccini/tosca"
)

func (self *Context) AddHierarchies() {
	self.Root.MergeHierarchies(make(tosca.HierarchyContext), self.HierarchiesWork)
}

func (self *Unit) MergeHierarchies(hierarchyContext tosca.HierarchyContext, work *ContextualWork) {
	context := self.GetContext()

	if promise, ok := work.Start(context); ok {
		defer promise.Release()

		for _, import_ := range self.Imports {
			import_.MergeHierarchies(hierarchyContext, work)
			context.Hierarchy.Merge(import_.GetContext().Hierarchy, hierarchyContext)
		}

		logHierarchies.Debugf("create: %s", context.URL.String())
		hierarchy := tosca.NewHierarchyFor(self.EntityPtr, hierarchyContext)
		context.Hierarchy.Merge(hierarchy, hierarchyContext)
		context.Hierarchy.AddTo(self.EntityPtr)
	}
}

// Print

func (self *Context) PrintHierarchies(indent int) {
	for _, import_ := range self.Units {
		context := import_.GetContext()
		if !context.Hierarchy.Empty() {
			terminal.PrintIndent(indent)
			fmt.Fprintf(terminal.Stdout, "%s\n", terminal.StyleValue(context.URL.String()))
			context.Hierarchy.Print(indent)
		}
	}
}
