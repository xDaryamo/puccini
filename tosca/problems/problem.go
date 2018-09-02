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
	URL     string
}

// fmt.Stringify interface
func (self Problem) String() string {
	return self.Message
}

//
// Problems
//

type Problems []Problem

func (self Problems) Empty() bool {
	return len(self) == 0
}

// fmt.Stringify interface
func (self Problems) String() string {
	var writer strings.Builder
	self.Write(&writer)
	return writer.String()
}

func (self Problems) Write(writer io.Writer) bool {
	length := len(self)
	if length > 0 {
		// Sort
		var problems = make(Problems, length)
		copy(problems, self)
		sort.Sort(problems)

		fmt.Fprintf(writer, "%s (%d)\n", format.ColorHeading("Problems"), length)
		var currentUrl string
		for _, problem := range problems {
			url_ := problem.URL
			if currentUrl != url_ {
				currentUrl = url_
				fmt.Fprint(writer, format.IndentString(1))
				if currentUrl != "" {
					fmt.Fprintf(writer, "%s\n", format.ColorValue(currentUrl))
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

func (self Problems) Print() bool {
	return self.Write(format.Stderr)
}

// sort.Interface

func (self Problems) Len() int {
	return len(self)
}

func (self Problems) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self Problems) Less(i, j int) bool {
	iProblem := self[i]
	jProblem := self[j]
	c := strings.Compare(iProblem.URL, jProblem.URL)
	if c == 0 {
		return strings.Compare(iProblem.Message, jProblem.Message) < 0
	}
	return c < 0
}
