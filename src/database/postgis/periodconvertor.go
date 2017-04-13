package postgis

import (
	"log"
	"strings"
	"time"

	"github.com/gost/now"
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

// ToPostgresPeriodFormat formats to Postgres format
func ToPostgresPeriodFormat(period [2]time.Time) string {
	start := ToIso8601(period[0])
	end := ToIso8601(period[1])
	return "[" + start + "," + end + "]"
}

// ToIso8601Period formats to Json Iso8601 period format
// sample input: ["2014-03-01 13:00:00+00","2015-05-11 15:30:00+00"]
// sample output:(2014-03-01T13:00:00Z/2015-05-11T15:30:00Z)
func ToIso8601Period(postgresPeriod string) string {
	log.Println(postgresPeriod)
	// log.Println(postgresPeriod)
	var iso8601Period = now.PostgresToIso8601Period(postgresPeriod)
	return iso8601Period
	// return "hoho" + postgresPeriod
}
