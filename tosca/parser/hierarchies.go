package parser

import (
	"fmt"

	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/tosca"
)

var hierarchiesWork = ContextsWork{Phase: "hierarchies"}

func (self *Context) AddHierarchies() {
	self.ServiceTemplate.MergeHierarchies(make(tosca.HierarchyContext))
}

func (self *Unit) MergeHierarchies(hierarchyContext tosca.HierarchyContext) {
	context := self.GetContext()

	if promise, ok := hierarchiesWork.Start(context); ok {
		defer promise.Release()

		for _, import_ := range self.Imports {
			import_.MergeHierarchies(hierarchyContext)
			context.Hierarchy.Merge(import_.GetContext().Hierarchy, hierarchyContext)
		}

		log.Infof("{hierarchies} create: %s", context.URL.String())
		hierarchy := tosca.NewHierarchy(self.EntityPtr, hierarchyContext)
		context.Hierarchy.Merge(hierarchy, hierarchyContext)
		context.Hierarchy.AddTo(self.EntityPtr)
	}
}

// Print

func (self *Context) PrintHierarchies(indent int) {
	for _, import_ := range self.Units {
		context := import_.GetContext()
		if len(context.Hierarchy.Children) > 0 {
			format.PrintIndent(indent)
			fmt.Fprintf(format.Stdout, "%s\n", format.ColorValue(context.URL.String()))
			context.Hierarchy.Print(indent)
		}
	}
}
