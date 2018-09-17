package parser

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/tliron/puccini/format"
	"github.com/tliron/puccini/tosca"
)

type Import struct {
	EntityPtr       interface{}
	Container       *Import
	Imports         Imports
	NameTransformer tosca.NameTransformer
	Locker          sync.Mutex
}

func NewImport(entityPtr interface{}, container *Import, nameTransformer tosca.NameTransformer) *Import {
	self := Import{
		EntityPtr:       entityPtr,
		Container:       container,
		NameTransformer: nameTransformer,
	}
	if container != nil {
		container.AddImport(&self)
	}
	return &self
}

func (self *Import) AddImport(import_ *Import) {
	self.Locker.Lock()
	self.Imports = append(self.Imports, import_)
	self.Locker.Unlock()
}

func (self *Import) GetContext() *tosca.Context {
	return tosca.GetContext(self.EntityPtr)
}

// Print

func (self *Import) PrintImports(indent int, treePrefix format.TreePrefix) {
	length := len(self.Imports)
	last := length - 1

	// Sort
	imports := make(Imports, length)
	copy(imports, self.Imports)
	sort.Sort(imports)

	for i, import_ := range imports {
		isLast := i == last
		import_.PrintNode(indent, treePrefix, isLast)
		import_.PrintImports(indent, append(treePrefix, isLast))
	}
}

func (self *Import) PrintNode(indent int, treePrefix format.TreePrefix, last bool) {
	treePrefix.Print(indent, last)
	fmt.Fprintf(format.Stdout, "%s\n", format.ColorValue(self.GetContext().URL.String()))
}

//
// Imports
//

type Imports []*Import

// sort.Interface

func (self Imports) Len() int {
	return len(self)
}

func (self Imports) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self Imports) Less(i, j int) bool {
	iName := self[i].GetContext().URL.String()
	jName := self[j].GetContext().URL.String()
	return strings.Compare(iName, jName) < 0
}
