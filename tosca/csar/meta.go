package csar

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Version 1.1 adds "Other-Definitions"
var MetaVersions = []Version{{1, 0}, {1, 1}}

var CsarVersions = []Version{{1, 1}}

//
// Meta
//

// See:
//   https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.3/TOSCA-Simple-Profile-YAML-v1.3.html#_Toc302251718
//   https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.2/TOSCA-Simple-Profile-YAML-v1.2.html#_Toc528072959
//   https://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/TOSCA-Simple-Profile-YAML-v1.1.html#_Toc489606742
//   https://docs.oasis-open.org/tosca/TOSCA/v1.0/TOSCA-v1.0.html#_Toc356403711

type Meta struct {
	Version          *Version
	CsarVersion      *Version
	CreatedBy        string
	EntryDefinitions string
	OtherDefinitions string
}

func ReadMeta(reader io.Reader) (*Meta, error) {
	map_, err := parseMeta(reader)
	if err != nil {
		return nil, err
	}

	if err = requireMeta(map_, "TOSCA-Meta-File-Version", "CSAR-Version", "Created-By"); err != nil {
		return nil, err
	}

	self := new(Meta)

	for name, value := range map_ {
		var err error
		switch name {
		case "TOSCA-Meta-File-Version":
			if self.Version, err = ParseVersion(value); err != nil {
				return nil, err
			}
			if !hasVersionMeta(*self.Version, MetaVersions) {
				return nil, fmt.Errorf("unsupported \"TOSCA-Meta-File-Version\" in TOSCA.meta: %s", self.Version.String())
			}

		case "CSAR-Version":
			if self.CsarVersion, err = ParseVersion(value); err != nil {
				return nil, err
			}
			if !hasVersionMeta(*self.CsarVersion, CsarVersions) {
				return nil, fmt.Errorf("unsupported \"CSAR-Version\" in TOSCA.meta: %s", self.CsarVersion.String())
			}

		case "Created-By":
			self.CreatedBy = value

		case "Entry-Definitions":
			self.EntryDefinitions = value

		case "Other-Definitions":
			// Added in TOSCA-Meta-File-Version 1.1
			self.OtherDefinitions = value
		}
	}

	return self, nil
}

func parseMeta(reader io.Reader) (map[string]string, error) {
	map_ := make(map[string]string)

	scanner := bufio.NewScanner(reader)

	var lineNumber uint
	var current string
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber += 1

		// Empty lines reset current name
		if line == "" {
			current = ""
			continue
		}

		// Lines beginning with space are appended to current name
		if strings.HasPrefix(line, " ") && (current != "") {
			map_[current] += line[1:]
			continue
		}

		split := strings.Split(line, ": ")
		if len(split) != 2 {
			return nil, fmt.Errorf("malformed line %d in TOSCA.meta: %s", lineNumber, line)
		}

		// New current name
		current = split[0]

		// Append to current
		map_[current] += split[1]
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return map_, nil
}

func requireMeta(data map[string]string, names ...string) error {
	for _, name := range names {
		if _, ok := data[name]; !ok {
			return fmt.Errorf("TOSCA.meta does not contain required %q", name)
		}
	}
	return nil
}

func hasVersionMeta(version Version, versions []Version) bool {
	for _, version_ := range versions {
		if version == version_ {
			return true
		}
	}
	return false
}
