package csar

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

var MetaVersion = Version{1, 0}
var CsarVersion = Version{1, 1}

//
// Meta
//

// See:
//  https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.2/os/TOSCA-Simple-Profile-YAML-v1.2-os.html#_Toc528072959
//  http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/os/TOSCA-Simple-Profile-YAML-v1.1-os.html#_Toc489606742
//  http://docs.oasis-open.org/tosca/TOSCA/v1.0/os/TOSCA-v1.0-os.html#_Toc356403713

type Meta struct {
	MetaVersion      *Version
	CsarVersion      *Version
	Creator          string
	EntryDefinitions string
}

func ReadMeta(reader io.Reader) (*Meta, error) {
	scanner := bufio.NewScanner(reader)

	data := make(map[string]string)

	var l uint
	var current string
	for scanner.Scan() {
		line := scanner.Text()
		l += 1

		// Empty lines reset current name
		if line == "" {
			current = ""
			continue
		}

		// Lines beginning with space are appended to current name
		if strings.HasPrefix(line, " ") && (current != "") {
			data[current] += line[1:]
			continue
		}

		split := strings.Split(line, ": ")
		if len(split) != 2 {
			return nil, fmt.Errorf("malformed line %d in TOSCA.meta: %s", l, line)
		}

		// New current name
		current = split[0]

		switch current {
		case "TOSCA-Meta-File-Version", "CSAR-Version", "Created-By", "Entry-Definitions":
			data[current] += split[1]
		default:
			return nil, fmt.Errorf("unsupported name in TOSCA.meta line %d: %s", l, current)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	self := &Meta{}

	for name, value := range data {
		var err error
		switch name {
		case "TOSCA-Meta-File-Version":
			if self.MetaVersion, err = ParseVersion(value); err != nil {
				return nil, err
			}
		case "CSAR-Version":
			if self.CsarVersion, err = ParseVersion(value); err != nil {
				return nil, err
			}
		case "Created-By":
			self.Creator = value
		case "Entry-Definitions":
			self.EntryDefinitions = value
		}
	}

	if err := require(data, "TOSCA-Meta-File-Version", "CSAR-Version", "Created-By"); err != nil {
		return nil, err
	}

	if *self.MetaVersion != MetaVersion {
		return nil, fmt.Errorf("unsupported \"TOSCA-Meta-File-Version\" in TOSCA.meta: %s", self.MetaVersion.String())
	}
	if *self.CsarVersion != CsarVersion {
		return nil, fmt.Errorf("unsupported \"CSAR-Version\" in TOSCA.meta: %s", self.CsarVersion.String())
	}

	return self, nil
}

func require(data map[string]string, names ...string) error {
	for _, name := range names {
		if _, ok := data[name]; !ok {
			return fmt.Errorf("TOSCA.meta does not contain required \"%s\"", name)
		}
	}
	return nil
}
