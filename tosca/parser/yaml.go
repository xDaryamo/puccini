package parser

import (
	"fmt"
	"strconv"
	"strings"
)

//
// YAMLError
//

type YAMLError struct {
	err error
}

func NewYAMLError(err error) *YAMLError {
	return &YAMLError{err}
}

// error interface
func (self *YAMLError) Error() string {
	return self.err.Error()
}

// fmt.Stringer interface
func (self *YAMLError) String() string {
	return self.err.Error()
}

// tosca.problems.Problematic interface
func (self *YAMLError) Problem() (string, string, int, int) {
	// Unfortunately, "gopkg.in/yaml.v3" just uses fmt.Errorf to create its errors,
	// so the only way we can extract line number information is by parsing the error string

	message := self.err.Error()
	if strings.HasPrefix(message, "yaml: line ") {
		suffix := message[11:]
		if colon := strings.Index(suffix, ": "); colon != -1 {
			line := suffix[:colon]
			if row, err := strconv.Atoi(line); err == nil {
				return fmt.Sprintf("malformed YAML, %s", suffix[colon+2:]), "", row, 0
			}
		}
	}

	return message, "", -1, -1
}
