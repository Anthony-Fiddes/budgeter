package period

// Period is an enum representing the lengths of time that budgeter allows
type Period int

const (
	Unknown Period = iota
	Day
	Week
	Month
)

func (p Period) Unknown() bool {
	if p <= 0 || p > Month {
		return true
	}
	return false
}

func (p Period) String() string {
	if p < 0 || p > Month {
		return "Unknown"
	}
	return [...]string{"Unknown", "Day", "Week", "Month"}[int(p)]
}

var periods = map[string]Period{
	Unknown.String(): Unknown,
	Day.String():     Day,
	Week.String():    Week,
	Month.String():   Month,
}

func Get(s string) Period {
	p, ok := periods[s]
	if !ok {
		return Unknown
	}
	return p
}
