package tosca_v2_0

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

var ScalarUnitTypeZeroes = map[ard.TypeName]ard.Value{
	"scalar-unit.size":      int(0),
	"scalar-unit.time":      float64(0.0),
	"scalar-unit.frequency": float64(0.0),
	"scalar-unit.bitrate":   float64(0.0),
}

//
// Scalar
//
// [TOSCA-v2.0] @ 9.1.2.2
//

type Scalar struct {
	CanonicalString string `json:"$string" yaml:"$string"`
	CanonicalNumber any    `json:"$number" yaml:"$number"` // float64 or int64
	OriginalString  string `json:"$originalString" yaml:"$originalString"`

	Scalar float64 `json:"scalar" yaml:"scalar"`
	Unit   string  `json:"unit" yaml:"unit"`

	// Fields exposed for JavaScript consumption
	DataTypeName  string             `json:"dataTypeName,omitempty"`
	BaseType      string             `json:"baseType,omitempty"`
	Units         map[string]float64 `json:"units,omitempty"`
	CanonicalUnit string             `json:"canonicalUnit,omitempty"`
	Prefixes      map[string]float64 `json:"prefixes,omitempty"`

	ScalarType    *DataType `traverse:"ignore" json:"-" yaml:"-"`
	canonicalUnit string
	dataType      string // "float" or "integer"
	units         map[string]float64
	prefixes      map[string]float64
}

func ReadScalar(context *parsing.Context, scalarType *DataType) *Scalar {
	self := &Scalar{
		ScalarType: scalarType,
	}

	if !context.ValidateType(ard.TypeString) {
		return self
	}

	self.OriginalString = *context.ReadString()

	if err := self.loadScalarTypeConfig(); err != nil {
		context.ReportValueMalformed("scalar", err.Error())
		return self
	}

	// Populate fields for JavaScript consumption
	if scalarType != nil {
		self.DataTypeName = scalarType.Name
		self.BaseType = self.dataType
		self.Units = self.units
		self.CanonicalUnit = self.canonicalUnit
		self.Prefixes = self.prefixes
	}

	// Parse format: number followed by whitespace followed by unit
	re := regexp.MustCompile(`^([+-]?[0-9]*\.?[0-9]+(?:[eE][+-]?[0-9]+)?)\s+(.+)$`)
	matches := re.FindStringSubmatch(self.OriginalString)
	if len(matches) != 3 {
		context.ReportValueMalformed("scalar", "format must be '<number> <unit>'")
		return self
	}

	// Parse scalar value
	var err error
	if self.Scalar, err = strconv.ParseFloat(matches[1], 64); err != nil {
		context.ReportValueMalformed("scalar", err.Error())
		return self
	}

	// Parse unit
	unitStr := matches[2]
	self.Unit = unitStr

	multiplier, err := self.getUnitMultiplier(unitStr)
	if err != nil {
		context.ReportValueMalformed("scalar", err.Error())
		return self
	}

	// Calculate canonical value
	canonicalValue := self.Scalar * multiplier

	if self.dataType == "integer" {
		self.CanonicalNumber = int64(math.Round(canonicalValue))
		self.CanonicalString = fmt.Sprintf("%d %s", self.CanonicalNumber, self.canonicalUnit)
	} else {
		self.CanonicalNumber = canonicalValue
		self.CanonicalString = fmt.Sprintf("%g %s", self.CanonicalNumber, self.canonicalUnit)
	}

	return self
}

func (self *Scalar) loadScalarTypeConfig() error {
	if self.ScalarType == nil {
		return errors.New("scalar type not set")
	}

	// Default data type is float
	self.dataType = "float"
	if self.ScalarType.DataTypeName != nil {
		self.dataType = *self.ScalarType.DataTypeName
	}

	// Load units from DataType
	self.units = make(map[string]float64)
	if self.ScalarType.Units != nil {
		for unitNameInterface, multiplierValue := range self.ScalarType.Units {
			if unitName, ok := unitNameInterface.(string); ok {
				if multiplier, ok := multiplierValue.(float64); ok {
					self.units[unitName] = multiplier
				} else if multiplierInt, ok := multiplierValue.(int64); ok {
					self.units[unitName] = float64(multiplierInt)
				} else if multiplierInt32, ok := multiplierValue.(int); ok {
					self.units[unitName] = float64(multiplierInt32)
				}
			}
		}
	}

	// Load prefixes from DataType
	self.prefixes = make(map[string]float64)
	if self.ScalarType.Prefixes != nil {
		for prefixNameInterface, multiplierValue := range self.ScalarType.Prefixes {
			if prefixName, ok := prefixNameInterface.(string); ok {
				if multiplier, ok := multiplierValue.(float64); ok {
					self.prefixes[prefixName] = multiplier
				} else if multiplierInt, ok := multiplierValue.(int64); ok {
					self.prefixes[prefixName] = float64(multiplierInt)
				} else if multiplierInt32, ok := multiplierValue.(int); ok {
					self.prefixes[prefixName] = float64(multiplierInt32)
				}
			}
		}
	}

	// Load canonical unit from DataType
	if self.ScalarType.CanonicalUnit != nil {
		self.canonicalUnit = *self.ScalarType.CanonicalUnit
	}

	// If no canonical unit specified, find the unit with multiplier 1
	if self.canonicalUnit == "" {
		for unit, multiplier := range self.units {
			if multiplier == 1.0 {
				self.canonicalUnit = unit
				break
			}
		}
	}

	return nil
}

func (self *Scalar) getUnitMultiplier(unitStr string) (float64, error) {
	// Try direct unit match first
	if multiplier, ok := self.units[unitStr]; ok {
		return multiplier, nil
	}

	// If prefixes are defined, try to match prefix + base unit
	if len(self.prefixes) > 0 {
		// Try to match prefix + base unit for each defined unit
		for unit, baseMultiplier := range self.units {
			// Try to match prefix + base unit
			for prefix, prefixMultiplier := range self.prefixes {
				if unitStr == prefix+unit {
					return baseMultiplier * prefixMultiplier, nil
				}
			}
		}
	}

	return 0, fmt.Errorf("unknown unit: %s", unitStr)
}

// ([fmt.Stringer] interface)
func (self *Scalar) String() string {
	return self.CanonicalString
}

func (self *Scalar) Compare(data any) (int, error) {
	if scalar, ok := data.(*Scalar); ok {
		if self.dataType == "integer" {
			return CompareInt64(self.CanonicalNumber.(int64), scalar.CanonicalNumber.(int64)), nil
		} else {
			return CompareFloat64(self.CanonicalNumber.(float64), scalar.CanonicalNumber.(float64)), nil
		}
	} else {
		return 0, errors.New("incompatible comparison")
	}
}

// ([parsing.Reader] signature) for generic scalar reading
func ReadScalarValue(context *parsing.Context) parsing.EntityPtr {
	context.ReportValueMalformed("scalar", "scalar value reading not yet fully implemented")
	return nil
}
