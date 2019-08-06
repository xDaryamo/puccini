package tosca_v1_3

import (
	"math"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
)

//
// Range
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.3.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.2.3
//

type Range struct {
	Lower uint64 `json:"lower" yaml:"lower"`
	Upper uint64 `json:"upper" yaml:"upper"`
}

// tosca.Reader signature
func ReadRange(context *tosca.Context) interface{} {
	var self Range

	if !context.ValidateType("list") {
		return self
	}

	list := context.Data.(ard.List)
	if len(list) != 2 {
		context.ReportValueMalformed("range", "list length not 2")
		return self
	}

	lowerContext := context.ListChild(0, list[0])
	lowerOk := false
	if lowerContext.ValidateType("integer") {
		lowerInt := *lowerContext.ReadInteger()
		if lowerInt < 0 {
			context.ReportValueMalformed("range", "lower bound negative")
		} else {
			self.Lower = uint64(lowerInt)
			lowerOk = true
		}
	}

	upperContext := context.ListChild(1, list[1])
	upperOk := false
	if upperContext.ValidateType("integer", "string") {
		if upperContext.Is("integer") {
			upperInt := *upperContext.ReadInteger()
			if upperInt < 0 {
				context.ReportValueMalformed("range", "upper bound negative")
			} else {
				self.Upper = uint64(upperInt)
				upperOk = true
			}
		} else if upperContext.Data.(string) == "UNBOUNDED" {
			self.Upper = math.MaxUint64
			upperOk = true
		} else {
			context.ReportValueMalformed("range", "upper bound string not UNBOUNDED")
		}
	}

	// The TOSCA spec says that the upper bound *must* be greater than the lower bound
	// but that makes no sense to us (it would mean than ranges must include at least two
	// numbers) so we are ignoring that and assuming that upper must be >= lower
	if upperOk && lowerOk && (self.Upper < self.Lower) {
		context.ReportValueMalformed("range", "upper bound lower than lower bound")
	}

	return self
}

func (self *Range) InRange(number uint64) bool {
	return (number >= self.Lower) && (number <= self.Upper)
}

//
// RangeEntity
//

type RangeEntity struct {
	*Entity `name:"range"`
	Range   Range
}

func NewRangeEntity(context *tosca.Context) *RangeEntity {
	return &RangeEntity{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadRangeEntity(context *tosca.Context) interface{} {
	self := NewRangeEntity(context)
	self.Range = ReadRange(context).(Range)
	return self
}
