package parser

import (
	contextpkg "context"
	"errors"

	"github.com/tliron/puccini/normal"
)

func (self *Context) Parse(context contextpkg.Context) (*normal.ServiceTemplate, error) {
	// Phase 1: Read
	ok := self.ReadRoot(context, self.URL, self.Bases, "")
	self.MergeProblems()
	problems := self.GetProblems()

	if !ok {
		return nil, errors.New("read error")
	}

	if !problems.Empty() {
		return nil, errors.New("read problems")
	}

	// Phase 2: Namespaces
	self.AddNamespaces()
	self.LookupNames()

	// Phase 3: Hierarchies
	self.AddHierarchies()

	// Phase 4: Inheritance
	self.Inherit(nil)

	self.SetInputs(self.Inputs)

	// Phase 5: Rendering
	self.Render()
	self.MergeProblems()
	if !problems.Empty() {
		return nil, errors.New("parsing problems")
	}

	// Phase 6: Normalization
	normalServiceTemplate, ok := self.Normalize()
	if !ok || !problems.Empty() {
		return nil, errors.New("normalization")
	}

	return normalServiceTemplate, nil
}
