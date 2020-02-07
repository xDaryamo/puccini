package tosca

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/tliron/puccini/common/terminal"
	"github.com/tliron/puccini/tosca/reflection"
)

//
// Hierarchical
//

type Hierarchical interface {
	GetParent() interface{}
}

// From Hierarchical interface
func GetParent(entityPtr interface{}) (interface{}, bool) {
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
	EntityPtr interface{}
	Parent    *Hierarchy
	Children  []*Hierarchy
}

// Keeps track of failed types
type HierarchyContext map[interface{}]bool

type HierarchyDescendants []interface{}

func NewHierarchy(entityPtr interface{}, hierarchyContext HierarchyContext) *Hierarchy {
	hierarchy := &Hierarchy{}
	reflection.Traverse(entityPtr, func(entityPtr interface{}) bool {
		if parentPtr, ok := GetParent(entityPtr); ok {
			hierarchy.Add(entityPtr, parentPtr, hierarchyContext, HierarchyDescendants{})
		}
		return true
	})
	return hierarchy
}

func (self *Hierarchy) GetCanonicalName() string {
	context := self.GetContext()
	canonicalNamespace := context.GetCanonicalNamespace()
	if canonicalNamespace != nil {
		return fmt.Sprintf("%s:%s", *canonicalNamespace, context.Name)
	} else {
		return context.Name
	}
}

func (self *Hierarchy) GetCanonicalNameFor(entityPtr interface{}) (string, bool) {
	if node, ok := self.Find(entityPtr); ok {
		return node.GetCanonicalName(), true
	} else {
		return "", false
	}
}

func (self *Hierarchy) Root() *Hierarchy {
	for self != nil {
		if self.Parent == nil {
			return self
		}
		self = self.Parent
	}
	panic("bad hierarchy")
}

func (self *Hierarchy) GetContext() *Context {
	return GetContext(self.EntityPtr)
}

func (self *Hierarchy) Find(entityPtr interface{}) (*Hierarchy, bool) {
	if entityPtr == nil {
		return nil, false
	}

	if self.EntityPtr == entityPtr {
		return self, true
	}

	for _, child := range self.Children {
		if node, ok := child.Find(entityPtr); ok {
			return node, true
		}
	}

	return nil, false
}

func (self *Hierarchy) Lookup(name string, type_ reflect.Type) (*Hierarchy, bool) {
	ourType := reflect.TypeOf(self.EntityPtr) == type_

	if ourType {
		if self.GetContext().Name == name {
			return self, true
		}
	}

	// Recurse only if in our type or in the root node
	if ourType || (self.Parent == nil) {
		for _, child := range self.Children {
			node, ok := child.Lookup(name, type_)
			if ok {
				return node, true
			}
		}
	}

	return nil, false
}

func (self *Hierarchy) IsCompatible(baseEntityPtr interface{}, entityPtr interface{}) bool {
	// Trivial case
	if baseEntityPtr == entityPtr {
		return true
	}

	if baseNode, ok := self.Find(baseEntityPtr); ok {
		_, ok = baseNode.Find(entityPtr)
		return ok
	}

	return false
}

func (self *Hierarchy) Add(entityPtr interface{}, parentPtr interface{}, hierarchyContext HierarchyContext, descendants HierarchyDescendants) (*Hierarchy, bool) {
	// Several imports may try to add the same entity to their hierarchies, so let's avoid multiple problem reports
	_, alreadyFailed := hierarchyContext[entityPtr]
	if alreadyFailed {
		return nil, false
	}

	root := self.Root()

	// Have we already added this entity?
	if found, ok := root.Find(entityPtr); ok {
		return found, true
	}

	node := &Hierarchy{EntityPtr: entityPtr}

	if parentPtr == nil {
		// We are a root node
		root.AddChild(node)
		return node, true
	}

	context := GetContext(entityPtr)

	// Check for inheritance loop
	for _, descendant := range descendants {
		if descendant == entityPtr {
			context.ReportInheritanceLoop(parentPtr)
			hierarchyContext[entityPtr] = true
			return nil, false
		}
	}

	grandparentPtr, ok := GetParent(parentPtr)
	if !ok {
		panic(fmt.Sprintf("parent is somehow of the wrong type (it doesn't have a \"parent\" tag): %s", context.Path))
	}

	// Make sure parent node has been added first (recursively)
	parentNode, ok := self.Add(parentPtr, grandparentPtr, hierarchyContext, append(descendants, entityPtr))
	if !ok {
		// Check if we already reported a failure (one report is enough)
		_, alreadyFailed := hierarchyContext[entityPtr]
		if !alreadyFailed {
			context.ReportTypeIncomplete(parentPtr)
			hierarchyContext[entityPtr] = true
		}
		return nil, false
	}

	// Add ourself to parent node
	parentNode.AddChild(node)

	return node, true
}

func (self *Hierarchy) Merge(node *Hierarchy, hierarchyContext HierarchyContext) {
	if node == nil {
		return
	}

	if node.EntityPtr != nil {
		if parentPtr, ok := GetParent(node.EntityPtr); ok {
			self.Add(node.EntityPtr, parentPtr, hierarchyContext, HierarchyDescendants{})
		}
	}

	for _, child := range node.Children {
		self.Merge(child, hierarchyContext)
	}
}

func (self *Hierarchy) AddChild(node *Hierarchy) {
	node.Parent = self
	self.Children = append(self.Children, node)
}

// Into "hierarchy" tags
func (self *Hierarchy) AddTo(entityPtr interface{}) {
	for _, field := range reflection.GetTaggedFields(entityPtr, "hierarchy") {
		type_ := field.Type()
		if reflection.IsSliceOfPtrToStruct(type_) {
			type_ = type_.Elem()
			self.AddTypeTo(field, type_)
		}
	}
}

func (self *Hierarchy) AddTypeTo(field reflect.Value, type_ reflect.Type) {
	if reflect.TypeOf(self.EntityPtr) == type_ {
		// Don't add if it's already there
		found := false
		length := field.Len()
		for i := 0; i < length; i++ {
			element := field.Index(i)
			if element.Interface() == self.EntityPtr {
				found = true
				break
			}
		}

		if !found {
			field.Set(reflect.Append(field, reflect.ValueOf(self.EntityPtr)))
		}
	}

	for _, child := range self.Children {
		child.AddTypeTo(field, type_)
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
	length := len(self.Children)
	last := length - 1

	// Sort
	hierarchy := Hierarchy{Children: make([]*Hierarchy, length)}
	copy(hierarchy.Children, self.Children)
	sort.Sort(hierarchy)

	for i, child := range hierarchy.Children {
		isLast := i == last
		child.PrintChild(indent, treePrefix, isLast)
		child.PrintChildren(indent, append(treePrefix, isLast))
	}
}

func (self *Hierarchy) PrintChild(indent int, treePrefix terminal.TreePrefix, last bool) {
	treePrefix.Print(indent, last)
	if self.EntityPtr != nil {
		fmt.Fprintf(terminal.Stdout, "%s\n", terminal.ColorTypeName(self.GetContext().Name))
	}
}

// sort.Interface

func (self Hierarchy) Len() int {
	return len(self.Children)
}

func (self Hierarchy) Swap(i, j int) {
	self.Children[i], self.Children[j] = self.Children[j], self.Children[i]
}

func (self Hierarchy) Less(i, j int) bool {
	iName := self.Children[i].GetContext().Name
	jName := self.Children[j].GetContext().Name
	return strings.Compare(iName, jName) < 0
}
