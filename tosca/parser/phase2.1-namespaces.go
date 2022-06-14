package parser

import (
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/puccini/tosca"
)

func (self *ServiceContext) AddNamespaces() {
	self.Context.lock.Lock()
	defer self.Context.lock.Unlock()

	self.Root.mergeNamespaces()
}

func (self *File) mergeNamespaces() {
	context := self.GetContext()

	self.importsLock.RLock()
	defer self.importsLock.RUnlock()

	for _, import_ := range self.Imports {
		import_.mergeNamespaces()
		context.Namespace.Merge(import_.GetContext().Namespace, import_.NameTransformer)
		context.ScriptletNamespace.Merge(import_.GetContext().ScriptletNamespace)
	}

	logNamespaces.Debugf("create: %s", context.URL.String())
	namespace := tosca.NewNamespaceFor(self.EntityPtr)
	context.Namespace.Merge(namespace, nil)
}

// Print

func (self *ServiceContext) PrintNamespaces(indent int) {
	self.filesLock.RLock()
	defer self.filesLock.RUnlock()

	childIndent := indent + 1
	for _, file := range self.Files {
		context := file.GetContext()
		if !context.Namespace.Empty() {
			terminal.PrintIndent(indent)
			terminal.Printf("%s\n", self.Stylist.Value(context.URL.String()))
			context.Namespace.Print(childIndent)
		}
	}
}
