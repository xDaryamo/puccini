package csar

import (
	"fmt"
	"strconv"
	"strings"
)

//
// Version
//

type Version struct {
	V1 uint8
	V2 uint8
}

func ParseVersion(value string) (*Version, error) {
	split := strings.Split(value, ".")
	if len(split) != 2 {
		return nil, fmt.Errorf("malformed version in TOSCA.meta: %s", value)
	}

	v1, err := parseDigit(value, split[0])
	if err != nil {
		return nil, err
	}

	v2, err := parseDigit(value, split[1])
	if err != nil {
		return nil, err
	}

	return &Version{v1, v2}, nil
}

func (self *Version) String() string {
	return fmt.Sprintf("%d.%d", self.V1, self.V2)
}

func parseDigit(value string, digit string) (uint8, error) {
	d, err := strconv.ParseUint(digit, 10, 8)

	// Seriously, the spec says this should be a single digit (?!)
	if (err != nil) || (d > 9) {
		return 0, fmt.Errorf("malformed version in TOSCA.meta: %s", value)
	}

	return uint8(d), nil
}
