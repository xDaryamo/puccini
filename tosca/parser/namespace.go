package parser

import (
	"fmt"

	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/tosca"
)

var namespacesWork = ContextsWork{Phase: "namespaces"}

func (self *Context) AddNamespaces() {
	self.ServiceTemplate.MergeNamespaces()
}

func (self *Unit) MergeNamespaces() {
	context := self.GetContext()

	if promise, ok := namespacesWork.Start(context); ok {
		defer promise.Release()

		for _, import_ := range self.Imports {
			import_.MergeNamespaces()
			context.Namespace.Merge(import_.GetContext().Namespace, import_.NameTransformer)
			context.ScriptletNamespace.Merge(import_.GetContext().ScriptletNamespace)
		}

		log.Infof("{namespaces} create: %s", context.URL.String())
		namespace := tosca.NewNamespace(self.EntityPtr)
		context.Namespace.Merge(namespace, nil)
	}
}

// Print

func (self *Context) PrintNamespaces(indent int) {
	childIndent := indent + 1
	for _, import_ := range self.Units {
		context := import_.GetContext()
		if len(context.Namespace) > 0 {
			format.PrintIndent(indent)
			fmt.Fprintf(format.Stdout, "%s\n", format.ColorValue(context.URL.String()))
			context.Namespace.Print(childIndent)
		}
	}
}
