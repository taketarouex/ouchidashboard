package enum

import "github.com/pkg/errors"

// Order desc or asc
type Order struct{ value string }

// Asc order
var Asc = Order{"ASC"}

// Desc order
var Desc = Order{"DESC"}

// LogType includes ouchi environment log types
type LogType struct{ value string }

// Temperature is a log type
var Temperature = LogType{"temperature"}

// Humidity same as above
var Humidity = LogType{"humidity"}

// Illumination same as above
var Illumination = LogType{"illumination"}

// Motion same as above
var Motion = LogType{"motion"}

func (t LogType) String() string {
	if t.value == "" {
		return "undefined"
	}
	return t.value
}

// ParseLogType parses string to LogType
func ParseLogType(target string) (LogType, error) {
	switch target {
	case Temperature.String():
		return Temperature, nil
	case Humidity.String():
		return Humidity, nil
	case Illumination.String():
		return Illumination, nil
	case Motion.String():
		return Motion, nil
	default:
		return LogType{}, errors.Errorf("invalid type: %s", target)
	}
}
