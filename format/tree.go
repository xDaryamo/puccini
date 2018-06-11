package format

import (
	"fmt"
)

//
// TreePrefix
//

type TreePrefix []bool

func (self TreePrefix) Print(indent int, last bool) {
	PrintIndent(indent)

	for _, element := range self {
		if element {
			fmt.Print("  ")
		} else {
			fmt.Print("│ ")
		}
	}

	if last {
		fmt.Print("└─")
	} else {
		fmt.Print("├─")
	}
}
