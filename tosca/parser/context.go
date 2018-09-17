package parser

import (
	"fmt"
	"sync"

	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/problems"
	"github.com/tliron/puccini/tosca/reflection"
)

type Context struct {
	ServiceTemplate *Import
	Problems        *problems.Problems
	Quirks          []string
	Imports         Imports
	Parsing         sync.Map
	WG              sync.WaitGroup
	Locker          sync.Mutex
}

func NewContext(quirks []string) Context {
	return Context{
		Problems: &problems.Problems{},
		Quirks:   quirks,
	}
}

func (self *Context) AddImport(import_ *Import) {
	self.Locker.Lock()
	self.Imports = append(self.Imports, import_)
	self.Locker.Unlock()
}

func (self *Context) AddImportFor(entityPtr interface{}, container *Import, nameTransformer tosca.NameTransformer) {
	import_ := NewImport(entityPtr, container, nameTransformer)
	if container == nil {
		self.AddImport(import_)
	}
	self.readImports(import_)
}

func (self *Context) Traverse(phase string, traverse reflection.Traverser) {
	done := make(EntitiesDone)
	t := func(entityPtr interface{}) bool {
		if done.IsDone(phase, entityPtr) {
			return false
		}
		return traverse(entityPtr)
	}

	reflection.Traverse(self.ServiceTemplate.EntityPtr, t)

	for _, forType := range self.ServiceTemplate.GetContext().Namespace {
		for _, entityPtr := range forType {
			reflection.Traverse(entityPtr, t)
		}
	}
}

// Print

func (self *Context) PrintImports(indent int) {
	format.PrintIndent(indent)
	fmt.Fprintf(format.Stdout, "%s\n", format.ColorValue(self.ServiceTemplate.GetContext().URL.String()))
	self.ServiceTemplate.PrintImports(indent, format.TreePrefix{})
}
