package tosca_v1_3

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/tliron/puccini/tosca"
)

// See the Timestamp Language-Independent Type for YAML Version 1.1 (Working Draft 2005-01-18)
// http://yaml.org/type/timestamp.html

var TimestampShortRE = regexp.MustCompile(`^(?P<year>[0-9][0-9][0-9][0-9])-(?P<month>[0-9][0-9])-(?P<day>[0-9][0-9])$`)

var TimestampLongRE = regexp.MustCompile(
	`^(?P<year>[0-9][0-9][0-9][0-9])-(?P<month>[0-9][0-9]?)-(?P<day>[0-9][0-9]?)` +
		`(?:[Tt]|[ \t]+)` +
		`(?P<hour>[0-9][0-9]?):(?P<minute>[0-9][0-9]):(?P<second>[0-9][0-9])(?:(?P<fraction>\.[0-9]*))?` +
		`(?:(?:[ \t]*)(?:Z|(?P<tzhour>[-+][0-9][0-9]?)(?::(?P<tzminute>[0-9][0-9]))?))?$`)

//
// Timestamp
//

type Timestamp struct {
	Number  int64  `json:"$number" yaml:"$number"`
	String_ string `json:"$string" yaml:"$string"`

	OriginalString string  `json:"originalString" yaml:"originalString"`
	Year           uint32  `json:"year" yaml:"year"`
	Month          uint32  `json:"month" yaml:"month"`
	Day            uint32  `json:"day" yaml:"day"`
	Hour           uint32  `json:"hour" yaml:"hour"`
	Minute         uint32  `json:"minute" yaml:"minute"`
	Second         uint32  `json:"second" yaml:"second"`
	Fraction       float64 `json:"fraction" yaml:"fraction"`
	TZSign         string  `json:"tzSign" yaml:"tzSign"`
	TZHour         uint32  `json:"tzHour" yaml:"tzHour"`
	TZMinute       uint32  `json:"tzMinute" yaml:"tzMinute"`
}

// tosca.Reader signature
func ReadTimestamp(context *tosca.Context) interface{} {
	var self Timestamp

	if !context.ValidateType("string") {
		return self
	}

	self.OriginalString = *context.ReadString()
	matches := TimestampShortRE.FindStringSubmatch(self.OriginalString)
	length := len(matches)
	if length == 0 {
		matches = TimestampLongRE.FindStringSubmatch(self.OriginalString)
		length = len(matches)
	}
	if length == 0 {
		context.ReportValueMalformed("timestamp", "")
		return self
	}

	if length > 1 {
		self.Year = parseTimestampUint(matches[1])
	}
	if length > 2 {
		self.Month = parseTimestampUint(matches[2])
		if (self.Month == 0) || (self.Month > 12) {
			context.ReportValueMalformed("timestamp", "invalid month")
			return self
		}
	}
	if length > 3 {
		self.Day = parseTimestampUint(matches[3])
		if self.Day > 31 {
			context.ReportValueMalformed("timestamp", "invalid day")
			return self
		}
	}
	if length > 4 {
		self.Hour = parseTimestampUint(matches[4])
		if self.Hour > 23 {
			context.ReportValueMalformed("timestamp", "invalid hour")
			return self
		}
	}
	if length > 5 {
		self.Minute = parseTimestampUint(matches[5])
		if self.Minute > 59 {
			context.ReportValueMalformed("timestamp", "invalid minute")
			return self
		}
	}
	if length > 6 {
		self.Second = parseTimestampUint(matches[6])
		if self.Second > 59 {
			context.ReportValueMalformed("timestamp", "invalid second")
			return self
		}
	}
	if (length > 7) && (matches[7] != "") {
		self.Fraction = parseTimestampFloat(matches[7])
	}
	if (length > 8) && (matches[8] != "") {
		v := matches[8]
		self.TZSign = v[:1]
		self.TZHour = parseTimestampUint(v[1:])

		// Real-world timezones go from -12 to +14:
		// https://en.wikipedia.org/wiki/List_of_UTC_time_offsets
		// But there's probably no benefit in validating that here, because the math of timezone
		// translation would "just work" even if the specified timezone doesn't exist in the real
		// world
	}
	if (length > 9) && (matches[9] != "") {
		self.TZMinute = parseTimestampUint(matches[9])
		if self.TZMinute > 59 {
			context.ReportValueMalformed("timestamp", "invalid TZ minute")
			return self
		}
	}

	// Canonical text
	var tz string
	if (self.TZHour == 0) && (self.TZMinute == 0) {
		tz = "Z"
	} else {
		tz = fmt.Sprintf("%s%02d:%02d", self.TZSign, self.TZHour, self.TZMinute)
	}
	fraction := fmt.Sprintf("%g", self.Fraction)[1:]
	self.String_ = fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d%s%s", self.Year, self.Month, self.Day, self.Hour, self.Minute, self.Second, fraction, tz)

	// Nanoseconds since Jan 1 1970 UTC
	self.Number = self.Time().UnixNano()

	return self
}

// fmt.Stringify interface
func (self *Timestamp) String() string {
	return self.String_
}

func (self *Timestamp) Compare(data interface{}) (int, error) {
	if timestamp, ok := data.(*Timestamp); ok {
		return CompareInt64(self.Number, timestamp.Number), nil
	}
	return 0, errors.New("incompatible comparison")
}

// Convert timezone to Go time.Location
func (self *Timestamp) Location() *time.Location {
	var factor int
	switch self.TZSign {
	case "-":
		factor = -1
	case "+":
		factor = 1
	}
	return time.FixedZone("", factor*int(self.TZHour)*3600+int(self.TZMinute)*60)
}

// Convert to Go time.Time
func (self *Timestamp) Time() time.Time {
	return time.Date(
		int(self.Year),
		time.Month(self.Month),
		int(self.Day),
		int(self.Hour),
		int(self.Minute),
		int(self.Second),
		int(self.Fraction*1000000000),
		self.Location(),
	)
}

func parseTimestampUint(value string) uint32 {
	if u, err := strconv.ParseUint(value, 10, 32); err == nil {
		return uint32(u)
	}
	panic("as long as the regexp does its job we should never get here")
}

func parseTimestampFloat(value string) float64 {
	if u, err := strconv.ParseFloat(value, 64); err == nil {
		return u
	}
	panic("as long as the regexp does its job we should never get here")
}
