package tosca_v1_3

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/tliron/puccini/tosca"
)

// Regexp

var ScalarUnitSizeRE = regexp.MustCompile(
	`^(?P<scalar>[0-9]*\.?[0-9]+(?:e[-+]?[0-9]+)?)\s*(?i)` +
		`(?P<unit>B|kB|KiB|MB|MiB|GB|GiB|TB|TiB)$`)

var ScalarUnitTimeRE = regexp.MustCompile(
	`^(?P<scalar>[0-9]*\.?[0-9]+(?:e[-+]?[0-9]+)?)\s*(?i)` +
		`(?P<unit>ns|us|ms|s|m|h|d)$`)

var ScalarUnitFrequencyRE = regexp.MustCompile(
	`^(?P<scalar>[0-9]*\.?[0-9]+(?:e[-+]?[0-9]+)?)\s*(?i)` +
		`(?P<unit>Hz|kHz|MHz|GHz)$`)

// Units

var ScalarUnitSizeSizes = ScalarUnitSizes{
	"B":   1,
	"kB":  1000,
	"KiB": 1024,
	"MB":  1000000,
	"MiB": 1048576,
	"GB":  1000000000,
	"GiB": 1073741824,
	"TB":  1000000000000,
	"TiB": 1099511627776,
}

var ScalarUnitTimeSizes = ScalarUnitSizes{
	"ns": 0.000000001,
	"us": 0.000001,
	"ms": 0.001,
	"s":  1,
	"m":  60,
	"h":  3600,
	"d":  86400,
}

var ScalarUnitFrequencySizes = ScalarUnitSizes{
	"Hz":  1,
	"kHz": 1000,
	"MHz": 1000000,
	"GHz": 1000000000,
}

//
// ScalarUnitSize
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.3.6.4
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.2.6.4
//

type ScalarUnitSize struct {
	CanonicalNumber uint64 `json:"$number" yaml:"$number"`
	CanonicalString string `json:"$string" yaml:"$string"`

	Scalar float64 `json:"scalar" yaml:"scalar"`
	Unit   string  `json:"unit" yaml:"unit"`

	OriginalString string `json:"originalString" yaml:"originalString"`
}

// tosca.Reader signature
func ReadScalarUnitSize(context *tosca.Context) interface{} {
	var self ScalarUnitSize

	originalString, scalar, unit, ok := parseScalarUnit(context, ScalarUnitSizeRE, "scalar-unit.size")
	if !ok {
		return self
	}

	normalUnit, size := ScalarUnitSizeSizes.Get(unit, context)

	self.OriginalString = originalString
	self.Scalar = scalar
	self.Unit = normalUnit
	self.CanonicalNumber = uint64(scalar * size)
	self.CanonicalString = fmt.Sprintf("%d B", self.CanonicalNumber)

	return self
}

// fmt.Stringify interface
func (self *ScalarUnitSize) String() string {
	if self.CanonicalNumber == 1 {
		return "1 byte"
	}
	return fmt.Sprintf("%d bytes", self.CanonicalNumber)
}

func (self *ScalarUnitSize) Compare(data interface{}) (int, error) {
	if scalarUnit, ok := data.(*ScalarUnitSize); ok {
		return CompareUint64(self.CanonicalNumber, scalarUnit.CanonicalNumber), nil
	}
	return 0, errors.New("incompatible comparison")
}

//
// ScalarUnitTime
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.3.6.5
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.2.6.5
//

type ScalarUnitTime struct {
	CanonicalNumber float64 `json:"$number" yaml:"$number"`
	CanonicalString string  `json:"$string" yaml:"$string"`

	Scalar float64 `json:"scalar" yaml:"scalar"`
	Unit   string  `json:"unit" yaml:"unit"`

	OriginalString string `json:"originalString" yaml:"originalString"`
}

// tosca.Reader signature
func ReadScalarUnitTime(context *tosca.Context) interface{} {
	var self ScalarUnitTime

	originalString, scalar, unit, ok := parseScalarUnit(context, ScalarUnitTimeRE, "scalar-unit.time")
	if !ok {
		return self
	}

	normalUnit, size := ScalarUnitTimeSizes.Get(unit, context)

	self.OriginalString = originalString
	self.Scalar = scalar
	self.Unit = normalUnit
	self.CanonicalNumber = scalar * size
	self.CanonicalString = fmt.Sprintf("%g S", self.CanonicalNumber)

	return self
}

// fmt.Stringify interface
func (self *ScalarUnitTime) String() string {
	if self.CanonicalNumber == 1.0 {
		return "1 second"
	}
	return fmt.Sprintf("%g seconds", self.CanonicalNumber)
}

func (self *ScalarUnitTime) Compare(data interface{}) (int, error) {
	if scalarUnit, ok := data.(*ScalarUnitTime); ok {
		return CompareFloat64(self.CanonicalNumber, scalarUnit.CanonicalNumber), nil
	}
	return 0, errors.New("incompatible comparison")
}

//
// ScalarUnitFrequency
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.3.6.6
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.2.6.6
//

type ScalarUnitFrequency struct {
	CanonicalNumber float64 `json:"$number" yaml:"$number"`
	CanonicalString string  `json:"$string" yaml:"$string"`

	Scalar float64 `json:"scalar" yaml:"scalar"`
	Unit   string  `json:"unit" yaml:"unit"`

	OriginalString string `json:"originalString" yaml:"originalString"`
}

// tosca.Reader signature
func ReadScalarUnitFrequency(context *tosca.Context) interface{} {
	var self ScalarUnitFrequency

	originalString, scalar, unit, ok := parseScalarUnit(context, ScalarUnitFrequencyRE, "scalar-unit.frequency")
	if !ok {
		return self
	}

	normalUnit, size := ScalarUnitFrequencySizes.Get(unit, context)

	self.OriginalString = originalString
	self.Scalar = scalar
	self.Unit = normalUnit
	self.CanonicalNumber = scalar * size
	self.CanonicalString = fmt.Sprintf("%g Hz", self.CanonicalNumber)

	return self
}

// fmt.Stringify interface
func (self *ScalarUnitFrequency) String() string {
	return fmt.Sprintf("%g Hz", self.CanonicalNumber)
}

func (self *ScalarUnitFrequency) Compare(data interface{}) (int, error) {
	if scalarUnit, ok := data.(*ScalarUnitFrequency); ok {
		return CompareFloat64(self.CanonicalNumber, scalarUnit.CanonicalNumber), nil
	}
	return 0, errors.New("incompatible comparison")
}

func parseScalarUnit(context *tosca.Context, re *regexp.Regexp, typeName string) (string, float64, string, bool) {
	if !context.ValidateType("string") {
		return "", 0, "", false
	}

	originalString := context.ReadString()
	matches := re.FindStringSubmatch(*originalString)
	if len(matches) != 3 {
		context.ReportValueMalformed(typeName, "")
		return "", 0, "", false
	}

	scalar, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		context.ReportValueMalformed(typeName, fmt.Sprintf("%s", err))
		return "", 0, "", false
	}

	return *originalString, scalar, matches[2], true
}

//
// ScalarUnitSizes
//

type ScalarUnitSizes map[string]float64

func (self ScalarUnitSizes) Get(unit string, context *tosca.Context) (string, float64) {
	for u, size := range self {
		if strings.EqualFold(u, unit) {
			return u, size
		}
	}
	panic("as long as the regexp does its job we should never get here")
}
