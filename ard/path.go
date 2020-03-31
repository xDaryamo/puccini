package ard

import (
	"fmt"
	"strings"
)

//
// PathElement
//

const (
	FieldPathType         = iota
	MapPathType           = iota
	ListPathType          = iota
	SequencedListPathType = iota
)

type PathElement struct {
	Type  int
	Value interface{} // string for FieldPathType and MapPathType, int for ListPathType
}

func NewFieldPathElement(name string) PathElement {
	return PathElement{FieldPathType, name}
}

func NewMapPathElement(name string) PathElement {
	return PathElement{MapPathType, name}
}

func NewListPathElement(index int) PathElement {
	return PathElement{ListPathType, index}
}

func NewSequencedListPathElement(index int) PathElement {
	return PathElement{SequencedListPathType, index}
}

//
// Path
//

type Path []PathElement

func (self Path) Append(element PathElement) Path {
	length := len(self)
	path := make(Path, length+1)
	copy(path, self)
	path[length] = element
	return path
}

func (self Path) AppendField(name string) Path {
	return self.Append(NewFieldPathElement(name))
}

func (self Path) AppendMap(name string) Path {
	return self.Append(NewMapPathElement(name))
}

func (self Path) AppendList(index int) Path {
	return self.Append(NewListPathElement(index))
}

func (self Path) AppendSequencedList(index int) Path {
	return self.Append(NewSequencedListPathElement(index))
}

// fmt.Stringer interface
func (self Path) String() string {
	var path string

	for _, element := range self {
		switch element.Type {
		case FieldPathType:
			value := escapeQuotes(element.Value.(string))
			if path == "" {
				path = value
			} else {
				path = fmt.Sprintf("%s.%s", path, value)
			}

		case MapPathType:
			value := escapeQuotes(element.Value.(string))
			path = fmt.Sprintf("%s[\"%s\"]", path, value)

		case ListPathType:
			value := element.Value.(int)
			path = fmt.Sprintf("%s[%d]", path, value)

		case SequencedListPathType:
			value := element.Value.(int)
			path = fmt.Sprintf("%s{%d}", path, value)
		}
	}

	return path
}

func escapeQuotes(s string) string {
	return strings.Replace(s, "\"", "\\\"", -1)
}
