package problems

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/tliron/puccini/format"
)

//
// Problem
//

type Problem struct {
	Message string
	Section string
}

// fmt.Stringify interface
func (self Problem) String() string {
	return self.Message
}

//
// ProblemSlice
//

type ProblemSlice []Problem

// sort.Interface

func (self ProblemSlice) Len() int {
	return len(self)
}

func (self ProblemSlice) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self ProblemSlice) Less(i, j int) bool {
	iProblem := self[i]
	jProblem := self[j]
	c := strings.Compare(iProblem.Section, jProblem.Section)
	if c == 0 {
		return strings.Compare(iProblem.Message, jProblem.Message) < 0
	}
	return c < 0
}

//
// Problems
//

type Problems struct {
	Problems ProblemSlice
}

func (self *Problems) Empty() bool {
	return len(self.Problems) == 0
}

// fmt.Stringify interface
func (self *Problems) String() string {
	var writer strings.Builder
	self.Write(&writer)
	return writer.String()
}

func (self *Problems) Write(writer io.Writer) bool {
	length := len(self.Problems)
	if length > 0 {
		// Sort
		problems := make(ProblemSlice, length)
		copy(problems, self.Problems)
		sort.Sort(problems)

		fmt.Fprintf(writer, "%s (%d)\n", format.ColorHeading("Problems"), length)
		var currentSection string
		for _, problem := range problems {
			section := problem.Section
			if currentSection != section {
				currentSection = section
				fmt.Fprint(writer, format.IndentString(1))
				if currentSection != "" {
					fmt.Fprintf(writer, "%s\n", format.ColorValue(currentSection))
				} else {
					fmt.Fprintf(writer, "General\n")
				}
			}
			fmt.Fprint(writer, format.IndentString(2))
			fmt.Fprintf(writer, "%s\n", problem)
		}
		return true
	}
	return false
}

// Print

func (self *Problems) Print() bool {
	return self.Write(format.Stderr)
}
