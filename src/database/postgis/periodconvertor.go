package postgis

import (
	"strings"
	"time"
)

// ToTime parses a time string to RCX3339Nano format
func ToTime(input string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, input)
}

// ToIso8601 formats a time object to ISO8601 string
func ToIso8601(t time.Time) string {
	return t.Format("2006-01-02T15:04:05.000Z")
}

// ParseTMPeriod parses a TM Period into time array (start,end)
func ParseTMPeriod(input string) [2]time.Time {
	parts := strings.Split(input, "/")
	startTime, _ := ToTime(parts[0])
	endTime, _ := ToTime(parts[1])
	var period [2]time.Time
	period[0] = startTime.UTC()
	period[1] = endTime.UTC()
	return period
}

// ToPostgresPeriodFormat formats to Posgres format
func ToPostgresPeriodFormat(period [2]time.Time) string {
	start := ToIso8601(period[0])
	end := ToIso8601(period[1])
	return "[" + start + "," + end + "]"
}
