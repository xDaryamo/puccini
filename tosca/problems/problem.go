package problems

import (
	"fmt"
	"io"
	"runtime"
	"sort"
	"strings"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/common/format"
	"github.com/tliron/puccini/common/terminal"
)

//
// Problem
//

type Problem struct {
	Message string `json:"message" yaml:"message"`
	Section string `json:"section" yaml:"section"`
	Row     int    `json:"row" yaml:"row"`
	Column  int    `json:"column" yaml:"column"`
	File    string `json:"file" yaml:"file"`
	Line    int    `json:"line" yaml:"line"`
}

func NewProblem(message string, section string, row int, column int, skip int) *Problem {
	self := Problem{
		Message: message,
		Section: section,
		Row:     row,
		Column:  column,
	}

	if _, file, line, ok := runtime.Caller(skip + 1); ok {
		self.File = file
		self.Line = line
	}

	return &self
}

// fmt.Stringer interface
func (self *Problem) String() string {
	r := ""
	if self.Row != -1 {
		r = fmt.Sprintf("@%d", self.Row)
		if self.Column != -1 {
			r = r + fmt.Sprintf(",%d", self.Column)
		}
		r = r + " "
	}
	r = r + self.Message
	return r
}

func (self *Problem) Equals(problem *Problem) bool {
	return (self.Message == problem.Message) && (self.Section == problem.Section) && (self.Row == problem.Row) && (self.Column == problem.Column)
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

// fmt.Stringer interface
func (self *Problems) String() string {
	return self.ToString(false)
}

func (self *Problems) ARD() (ard.Map, error) {
	if s, err := format.EncodeYAML(self, " ", false); err == nil {
		map_, _, err := ard.ReadYAML(strings.NewReader(s), false)
		return map_, err
	} else {
		return nil, err
	}
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
