package problems

import (
	"fmt"
	"io"
	"runtime"
	"sort"
	"strings"

	"github.com/tliron/puccini/common/terminal"
)

//
// Problem
//

type Problem struct {
	Message string
	Section string
	File    string
	Line    int
}

func NewProblem(message string, section string, skip int) *Problem {
	self := Problem{
		Message: message,
		Section: section,
	}

	if _, file, line, ok := runtime.Caller(skip + 1); ok {
		self.File = file
		self.Line = line
	}

	return &self
}

// fmt.Stringify interface
func (self *Problem) String() string {
	return self.Message
}

func (self *Problem) Equals(problem *Problem) bool {
	// TODO: compare File and Line?
	return (self.Message == problem.Message) && (self.Section == problem.Section)
}

//
// ProblemSlice
//

type ProblemSlice []*Problem

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

func (self *Problems) Append(problem *Problem) bool {
	// Avoid duplicates
	for _, problem_ := range self.Problems {
		if problem.Equals(problem_) {
			return false
		}
	}

	self.Problems = append(self.Problems, problem)
	return true
}

func (self *Problems) Merge(problems *Problems) bool {
	if self == problems {
		// Merging into self
		return false
	}

	merged := false
	for _, problem := range problems.Problems {
		if self.Append(problem) {
			merged = true
		}
	}

	return merged
}

func (self *Problems) ToString(locate bool) string {
	var writer strings.Builder
	self.Write(&writer, false, locate)
	return writer.String()
}

// fmt.Stringify interface
func (self *Problems) String() string {
	return self.ToString(false)
}

func (self *Problems) Write(writer io.Writer, pretty bool, locate bool) bool {
	length := len(self.Problems)
	if length > 0 {
		// Sort
		problems := make(ProblemSlice, length)
		copy(problems, self.Problems)
		sort.Sort(problems)

		if pretty {
			fmt.Fprintf(writer, "%s (%d)\n", terminal.ColorHeading("Problems"), length)
		} else {
			fmt.Fprintf(writer, "%s (%d)\n", "Problems", length)
		}

		var currentSection string
		for _, problem := range problems {
			section := problem.Section
			if currentSection != section {
				currentSection = section
				fmt.Fprint(writer, terminal.IndentString(1))
				if currentSection != "" {
					if pretty {
						fmt.Fprintf(writer, "%s\n", terminal.ColorValue(currentSection))
					} else {
						fmt.Fprintf(writer, "%s\n", currentSection)
					}
				} else {
					fmt.Fprintf(writer, "General\n")
				}
			}

			fmt.Fprint(writer, terminal.IndentString(2))
			fmt.Fprintf(writer, "%s\n", problem)

			if locate && (problem.File != "") {
				fmt.Fprint(writer, terminal.IndentString(2))
				fmt.Fprintf(writer, "└─%s:%d\n", problem.File, problem.Line)
			}
		}
		return true
	}
	return false
}

// Print

func (self *Problems) Print(locate bool) bool {
	return self.Write(terminal.Stderr, true, locate)
}
