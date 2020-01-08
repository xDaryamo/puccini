package parser

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/tosca"
)

type Unit struct {
	EntityPtr       interface{}
	Container       *Unit
	Imports         Units
	NameTransformer tosca.NameTransformer
	Locker          sync.Mutex
}

func NewUnit(entityPtr interface{}, container *Unit, nameTransformer tosca.NameTransformer) *Unit {
	self := Unit{
		EntityPtr:       entityPtr,
		Container:       container,
		NameTransformer: nameTransformer,
	}
	if container != nil {
		container.AddImport(&self)
	}
	return &self
}

func (self *Unit) AddImport(import_ *Unit) {
	self.Locker.Lock()
	self.Imports = append(self.Imports, import_)
	self.Locker.Unlock()
}

func (self *Unit) GetContext() *tosca.Context {
	return tosca.GetContext(self.EntityPtr)
}

// Print

func (self *Unit) PrintImports(indent int, treePrefix format.TreePrefix) {
	length := len(self.Imports)
	last := length - 1

	// Sort
	imports := make(Units, length)
	copy(imports, self.Imports)
	sort.Sort(imports)

	for i, unit := range imports {
		isLast := i == last
		unit.PrintNode(indent, treePrefix, isLast)
		unit.PrintImports(indent, append(treePrefix, isLast))
	}
}

func (self *Unit) PrintNode(indent int, treePrefix format.TreePrefix, last bool) {
	treePrefix.Print(indent, last)
	fmt.Fprintf(format.Stdout, "%s\n", format.ColorValue(self.GetContext().URL.String()))
}

//
// Units
//

type Units []*Unit

// sort.Interface

func (self Units) Len() int {
	return len(self)
}

func (self Units) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self Units) Less(i, j int) bool {
	iName := self[i].GetContext().URL.String()
	jName := self[j].GetContext().URL.String()
	return strings.Compare(iName, jName) < 0
}
