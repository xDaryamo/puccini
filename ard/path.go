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
	Value interface{}
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
			if path == "" {
				path = fmt.Sprintf("%s", element.Value.(string))
			} else {
				path = fmt.Sprintf("%s.%s", path, element.Value.(string))
			}

		case MapPathType:
			path = fmt.Sprintf("%s[\"%s\"]", path, strings.Replace(element.Value.(string), "\"", "\\\"", -1))

		case ListPathType:
			path = fmt.Sprintf("%s[%d]", path, element.Value.(int))
		}
	}

	return path
}
