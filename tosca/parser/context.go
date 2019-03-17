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
	ServiceTemplate *Unit
	Problems        problems.Problems
	Quirks          []string
	Units           Units
	Parsing         sync.Map
	WG              sync.WaitGroup
	Locker          sync.Mutex
}

func NewContext(quirks []string) Context {
	return Context{
		Quirks: quirks,
	}
}

func (self *Context) AddUnit(unit *Unit) {
	self.Locker.Lock()
	self.Units = append(self.Units, unit)
	self.Locker.Unlock()
}

func (self *Context) AddUnitFor(entityPtr interface{}, container *Unit, nameTransformer tosca.NameTransformer) {
	unit := NewUnit(entityPtr, container, nameTransformer)
	if container == nil {
		// It's a root unit, so it won't be added later
		self.AddUnit(unit)
	}
	self.goReadImports(unit)
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
