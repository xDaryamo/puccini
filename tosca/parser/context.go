package parser

import (
	"sync"

	"github.com/tliron/kutil/problems"
	"github.com/tliron/kutil/reflection"
	"github.com/tliron/kutil/terminal"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars"
)

//
// Context
//

type Context struct {
	readCache          sync.Map // entityPtr or Promise
	lookupFieldsWork   reflection.EntityWork
	addHierarchyWork   reflection.EntityWork
	getInheritTaskWork reflection.EntityWork
	renderWork         reflection.EntityWork
	lock               util.RWLocker
}

func NewContext() *Context {
	return &Context{
		lookupFieldsWork:   make(reflection.EntityWork),
		addHierarchyWork:   make(reflection.EntityWork),
		getInheritTaskWork: make(reflection.EntityWork),
		renderWork:         make(reflection.EntityWork),
		lock:               util.NewDefaultRWLocker(),
	}
}

//
// ServiceContext
//

type ServiceContext struct {
	Context *Context
	Root    *File
	Stylist *terminal.Stylist
	Quirks  tosca.Quirks
	Files   Files

	readWork  sync.WaitGroup
	filesLock util.RWLocker
}

func (self *Context) NewServiceContext(stylist *terminal.Stylist, quirks tosca.Quirks) *ServiceContext {
	return &ServiceContext{
		Context:   self,
		Stylist:   stylist,
		Quirks:    quirks,
		filesLock: util.NewDebugRWLocker(),
	}
}

func (self *ServiceContext) GetProblems() *problems.Problems {
	return self.Root.GetContext().Problems
}

func (self *ServiceContext) MergeProblems() {
	self.filesLock.RLock()
	defer self.filesLock.RUnlock()

	// Note: This could happen many times, but because problems are de-duped, everything is OK :)
	for _, file := range self.Files {
		self.GetProblems().Merge(file.GetContext().Problems)
	}
}

func (self *ServiceContext) AddFile(file *File) {
	self.filesLock.Lock()
	defer self.filesLock.Unlock()

	self.Files = append(self.Files, file)
}

func (self *ServiceContext) AddImportFile(entityPtr tosca.EntityPtr, container *File, nameTransformer tosca.NameTransformer) *File {
	file := NewFile(entityPtr, container, nameTransformer)

	if container != nil {
		containerContext := container.GetContext()
		if !containerContext.HasQuirk(tosca.QuirkImportsVersionPermissive) {
			fileContext := file.GetContext()
			if !grammars.CompatibleGrammars(containerContext, fileContext) {
				containerContext.ReportImportIncompatible(fileContext.URL)
				return file
			}
		}
	}

	self.AddFile(file)

	self.goReadImports(file)

	return file
}

// Print

func (self *ServiceContext) PrintImports(indent int) {
	terminal.PrintIndent(indent)
	terminal.Printf("%s\n", terminal.DefaultStylist.Value(self.Root.GetContext().URL.String()))
	self.Root.PrintImports(indent, terminal.TreePrefix{})
}
