package tosca

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/tliron/kutil/reflection"
	"github.com/tliron/kutil/terminal"
)

//
// Hierarchical
//

type Hierarchical interface {
	GetParent() EntityPtr
}

// From Hierarchical interface
func GetParent(entityPtr EntityPtr) (EntityPtr, bool) {
	if hierarchical, ok := entityPtr.(Hierarchical); ok {
		parent := hierarchical.GetParent()
		parentValue := reflect.ValueOf(parent)
		if !parentValue.IsNil() {
			return parent, true
		} else {
			return nil, true
		}
	}
	return nil, false
}

//
// Hierarchy
//

type Hierarchy struct {
	entityPtr EntityPtr
	parent    *Hierarchy
	children  []*Hierarchy
	lock      *sync.RWMutex // shared recursively with children
}

// Keeps track of failed types
type HierarchyContext map[EntityPtr]bool

type HierarchyDescendants []EntityPtr

func NewHierarchy() *Hierarchy {
	return &Hierarchy{
		lock: new(sync.RWMutex),
	}
}

func NewHierarchyFor(entityPtr EntityPtr, hierarchyContext HierarchyContext) *Hierarchy {
	self := NewHierarchy()

	reflection.Traverse(entityPtr, func(entityPtr EntityPtr) bool {
		if parentPtr, ok := GetParent(entityPtr); ok {
			self.add(entityPtr, parentPtr, hierarchyContext, HierarchyDescendants{})
		}
		return true
	})

	return self
}

func (self *Hierarchy) Empty() bool {
	self.lock.RLock()
	defer self.lock.RUnlock()

	return len(self.children) == 0
}

func (self *Hierarchy) Root() *Hierarchy {
	self.lock.RLock()
	defer self.lock.RUnlock()

	return self.root()
}

func (self *Hierarchy) root() *Hierarchy {
	for self != nil {
		if self.parent == nil {
			return self
		}
		self = self.parent
	}

	panic("bad hierarchy")
}

func (self *Hierarchy) GetContext() *Context {
	self.lock.RLock()
	defer self.lock.RUnlock()

	return self.getContext()
}

func (self *Hierarchy) getContext() *Context {
	return GetContext(self.entityPtr)
}

func (self *Hierarchy) Range(f func(EntityPtr, EntityPtr) bool) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	hierarchy_ := self
	for (hierarchy_ != nil) && (hierarchy_.entityPtr != nil) {
		parent := hierarchy_.parent

		if parent != nil {
			if !f(hierarchy_.entityPtr, parent.entityPtr) {
				return
			}
		} else if !f(hierarchy_.entityPtr, nil) {
			return
		}

		hierarchy_ = parent
	}
}

func (self *Hierarchy) Find(entityPtr EntityPtr) (*Hierarchy, bool) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	return self.find(entityPtr)
}

func (self *Hierarchy) find(entityPtr EntityPtr) (*Hierarchy, bool) {
	if entityPtr == nil {
		return nil, false
	}

	if self.entityPtr == entityPtr {
		return self, true
	}

	for _, child := range self.children {
		if found, ok := child.find(entityPtr); ok {
			return found, true
		}
	}

	return nil, false
}

func (self *Hierarchy) Lookup(name string, type_ reflect.Type) (*Hierarchy, bool) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	return self.lookup(name, type_)
}

func (self *Hierarchy) lookup(name string, type_ reflect.Type) (*Hierarchy, bool) {
	ourType := reflect.TypeOf(self.entityPtr) == type_

	if ourType {
		if self.getContext().Name == name {
			return self, true
		}
	}

	// Recurse only if in our type or in the root node
	if ourType || (self.parent == nil) {
		for _, child := range self.children {
			if found, ok := child.lookup(name, type_); ok {
				return found, true
			}
		}
	}

	return nil, false
}

func (self *Hierarchy) IsCompatible(baseEntityPtr EntityPtr, entityPtr EntityPtr) bool {
	// Trivial case
	if baseEntityPtr == entityPtr {
		return true
	}

	self.lock.RLock()
	defer self.lock.RUnlock()

	if baseNode, ok := self.find(baseEntityPtr); ok {
		_, ok = baseNode.find(entityPtr)
		return ok
	}

	return false
}

func (self *Hierarchy) add(entityPtr EntityPtr, parentEntityPtr EntityPtr, hierarchyContext HierarchyContext, descendants HierarchyDescendants) (*Hierarchy, bool) {
	// Several imports may try to add the same entity to their hierarchies, so let's avoid multiple problem reports
	_, alreadyFailed := hierarchyContext[entityPtr]
	if alreadyFailed {
		return nil, false
	}

	root := self.root()

	// Have we already added this entity?
	if found, ok := root.find(entityPtr); ok {
		return found, true
	}

	child := &Hierarchy{
		entityPtr: entityPtr,
		lock:      self.lock,
	}

	if parentEntityPtr == nil {
		// We are a root node
		root.addChild(child)
		return child, true
	}

	context := GetContext(entityPtr)

	// Check for inheritance loop
	for _, descendant := range descendants {
		if descendant == entityPtr {
			context.ReportInheritanceLoop(parentEntityPtr)
			hierarchyContext[entityPtr] = true
			return nil, false
		}
	}

	grandparentEntityPtr, ok := GetParent(parentEntityPtr)
	if !ok {
		panic(fmt.Sprintf("parent is somehow of the wrong type (it doesn't have a \"parent\" tag): %s", context.Path))
	}

	// Make sure parent node has been added first (recursively)
	parentNode, ok := self.add(parentEntityPtr, grandparentEntityPtr, hierarchyContext, append(descendants, entityPtr))
	if !ok {
		// Check if we already reported a failure (one report is enough)
		_, alreadyFailed := hierarchyContext[entityPtr]
		if !alreadyFailed {
			context.ReportTypeIncomplete(parentEntityPtr)
			hierarchyContext[entityPtr] = true
		}
		return nil, false
	}

	// Add ourself to parent node
	parentNode.addChild(child)

	return child, true
}

func (self *Hierarchy) Merge(hierarchy *Hierarchy, hierarchyContext HierarchyContext) {
	if (hierarchy == nil) || (self == hierarchy) {
		return
	}

	self.lock.Lock()
	defer self.lock.Unlock()

	if self.lock != hierarchy.lock {
		hierarchy.lock.RLock()
		defer hierarchy.lock.RUnlock()
	}

	self.merge(hierarchy, hierarchyContext)
}

func (self *Hierarchy) merge(hierarchy *Hierarchy, hierarchyContext HierarchyContext) {
	if hierarchy.entityPtr != nil {
		if parentPtr, ok := GetParent(hierarchy.entityPtr); ok {
			self.add(hierarchy.entityPtr, parentPtr, hierarchyContext, HierarchyDescendants{})
		}
	}

	for _, child := range hierarchy.children {
		self.merge(child, hierarchyContext)
	}
}

func (self *Hierarchy) addChild(hierarchy *Hierarchy) {
	hierarchy.parent = self
	self.children = append(self.children, hierarchy)
}

// Into "hierarchy" tags
func (self *Hierarchy) AddTo(entityPtr EntityPtr) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	for _, field := range reflection.GetTaggedFields(entityPtr, "hierarchy") {
		type_ := field.Type()
		if reflection.IsSliceOfPtrToStruct(type_) {
			type_ = type_.Elem()
			self.addTypeTo(field, type_)
		}
	}
}

func (self *Hierarchy) addTypeTo(field reflect.Value, type_ reflect.Type) {
	if reflect.TypeOf(self.entityPtr) == type_ {
		// Don't add if it's already there
		found := false
		length := field.Len()
		for i := 0; i < length; i++ {
			element := field.Index(i)
			if element.Interface() == self.entityPtr {
				found = true
				break
			}
		}

		if !found {
			field.Set(reflect.Append(field, reflect.ValueOf(self.entityPtr)))
		}
	}

	for _, child := range self.children {
		child.addTypeTo(field, type_)
	}
}

// Print

// Note that the same name could be printed out twice in the hierarchy, even under the same
// parent! That's because we are printing the local name of the type, and types imported from
// other files can have the same name (though you would need a namespace_prefix to avoid a
// namespace error)
func (self *Hierarchy) Print(indent int) {
	self.PrintChildren(indent, terminal.TreePrefix{})
}

func (self *Hierarchy) PrintChildren(indent int, treePrefix terminal.TreePrefix) {
	length := len(self.children)
	last := length - 1

	// Sort
	self.lock.RLock()
	hierarchy := Hierarchy{
		children: make([]*Hierarchy, length),
		lock:     self.lock,
	}
	copy(hierarchy.children, self.children)
	self.lock.RUnlock()
	sort.Sort(hierarchy)

	for i, child := range hierarchy.children {
		isLast := i == last
		child.PrintChild(indent, treePrefix, isLast)
		child.PrintChildren(indent, append(treePrefix, isLast))
	}
}

func (self *Hierarchy) PrintChild(indent int, treePrefix terminal.TreePrefix, last bool) {
	treePrefix.Print(indent, last)
	if self.entityPtr != nil {
		fmt.Fprintf(terminal.Stdout, "%s\n", terminal.StyleTypeName(self.GetContext().Name))
	}
}

// sort.Interface

func (self Hierarchy) Len() int {
	return len(self.children)
}

func (self Hierarchy) Swap(i, j int) {
	self.children[i], self.children[j] = self.children[j], self.children[i]
}

func (self Hierarchy) Less(i, j int) bool {
	iName := self.children[i].GetContext().Name
	jName := self.children[j].GetContext().Name
	return strings.Compare(iName, jName) < 0
}
