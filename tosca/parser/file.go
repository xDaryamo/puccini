package parser

import (
	"sort"
	"strings"

	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca"
)

//
// NoEntity
//

type NoEntity struct {
	Context *tosca.Context
}

func NewNoEntity(toscaContext *tosca.Context) *NoEntity {
	return &NoEntity{toscaContext}
}

// tosca.Contextual interface
func (self *NoEntity) GetContext() *tosca.Context {
	return self.Context
}

//
// File
//

type File struct {
	EntityPtr       tosca.EntityPtr
	Container       *File
	Imports         Files
	NameTransformer tosca.NameTransformer

	importsLock util.RWLocker
}

func NewFileNoEntity(toscaContext *tosca.Context, container *File, nameTransformer tosca.NameTransformer) *File {
	return NewFile(NewNoEntity(toscaContext), container, nameTransformer)
}

func NewFile(entityPtr tosca.EntityPtr, container *File, nameTransformer tosca.NameTransformer) *File {
	self := File{
		EntityPtr:       entityPtr,
		Container:       container,
		NameTransformer: nameTransformer,
		importsLock:     util.NewDebugRWLocker(),
	}

	if container != nil {
		container.AddImport(&self)
	}

	return &self
}

func (self *File) AddImport(import_ *File) {
	self.importsLock.Lock()
	defer self.importsLock.Unlock()

	self.Imports = append(self.Imports, import_)
}

func (self *File) GetContext() *tosca.Context {
	return tosca.GetContext(self.EntityPtr)
}

// Print

func (self *File) PrintImports(indent int, treePrefix terminal.TreePrefix) {
	self.importsLock.RLock()
	length := len(self.Imports)
	imports := make(Files, length)
	copy(imports, self.Imports)
	self.importsLock.RUnlock()

	last := length - 1

	// Sort
	sort.Sort(imports)

	for i, file := range imports {
		isLast := i == last
		file.PrintNode(indent, treePrefix, isLast)
		file.PrintImports(indent, append(treePrefix, isLast))
	}
}

func (self *File) PrintNode(indent int, treePrefix terminal.TreePrefix, last bool) {
	treePrefix.Print(indent, last)
	terminal.Printf("%s\n", terminal.DefaultStylist.Value(self.GetContext().URL.String()))
}

//
// Files
//

type Files []*File

// sort.Interface

func (self Files) Len() int {
	return len(self)
}

func (self Files) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self Files) Less(i, j int) bool {
	iName := self[i].GetContext().URL.String()
	jName := self[j].GetContext().URL.String()
	return strings.Compare(iName, jName) < 0
}
