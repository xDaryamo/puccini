package csar

import (
	"bufio"
	"bytes"
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
	Version          *Version `yaml:"version" json:"version"`
	CsarVersion      *Version `yaml:"csarVersion" json:"csarVersion"`
	CreatedBy        string   `yaml:"createdBy" json:"createdBy"`
	EntryDefinitions string   `yaml:"entryDefinitions" json:"entryDefinitions"`
	OtherDefinitions []string `yaml:"otherDefinitions" json:"otherDefinitions"`
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
			if self.OtherDefinitions, err = parseStringList(value); err != nil {
				return nil, err
			}
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

func parseStringList(value string) ([]string, error) {
	// Note: The TOSCA specification does not mention the possibility of escaping quotation marks,
	// but it should obviously be supported. So we are reserving backslashes for escaping.
	var entries []string
	var entry bytes.Buffer
	escaped := false
	spaced := false
	quoted := false

	for _, rune_ := range value {
		if escaped {
			escaped = false
			entry.WriteRune(rune_)
			continue
		}

		switch rune_ {
		case '\\':
			spaced = false
			escaped = true

		case ' ':
			// The spec says "a blank space", so we will treat multiple spaces as an error
			if spaced {
				return nil, fmt.Errorf("malformed string list, separator must be single space: %s", value)
			} else if quoted {
				entry.WriteRune(rune_)
			} else {
				spaced = true
				if entry_ := entry.String(); len(entry_) > 0 {
					entries = append(entries, entry_)
					entry.Reset()
				}
			}

		case '"':
			spaced = false
			if quoted {
				// End quote
				quoted = false
				if entry_ := entry.String(); len(entry_) > 0 {
					entries = append(entries, entry_)
					entry.Reset()
				}
			} else {
				// Start quote
				quoted = true
			}

		default:
			spaced = false
			entry.WriteRune(rune_)
		}
	}

	if escaped {
		return nil, fmt.Errorf("malformed string list, ends with a backslash: %s", value)
	}

	if spaced {
		return nil, fmt.Errorf("malformed string list, ends with a space: %s", value)
	}

	if quoted {
		return nil, fmt.Errorf("malformed string list, did not close quotation: %s", value)
	}

	// Last entry
	if entry_ := entry.String(); len(entry_) > 0 {
		entries = append(entries, entry_)
	}

	return entries, nil
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
