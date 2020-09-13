package enum

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
