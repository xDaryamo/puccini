package ard

import (
	"fmt"
	"strings"
)

//
// PathElement
//

const (
	FieldPathType = iota
	MapPathType   = iota
	ListPathType  = iota
)

type PathElement struct {
	Type  int
	Value interface{} // int for ListPathType, string for FieldPathType and MapPathType
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

//
// Path
//

type Path []PathElement

func (self Path) String() string {
	var path string

	for _, element := range self {
		switch element.Type {
		case FieldPathType:
			value := element.Value.(string)
			if path == "" {
				path = value
			} else {
				path = fmt.Sprintf("%s.%s", path, value)
			}

		case MapPathType:
			value := element.Value.(string)
			path = fmt.Sprintf("%s[\"%s\"]", path, escapeQuotes(value))

		case ListPathType:
			value := element.Value.(int)
			path = fmt.Sprintf("%s[%d]", path, value)
		}
	}

	return path
}

func escapeQuotes(s string) string {
	return strings.Replace(s, "\"", "\\\"", -1)
}
