package postgis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIso8601(t *testing.T) {
	// arrange
	var testdate = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	// act
	isodate := ToIso8601(testdate)
	// assert
	assert.Equal(t, "2009-11-10T23:00:00.000Z", isodate, "Iso8601 error")
}

func TestToTime(t *testing.T) {
	// arrange
	var testtime = "2014-03-01T13:00:00Z"
	var expectedtime = time.Date(2014, time.March, 1, 13, 0, 0, 0, time.UTC)
	// act
	gotime, _ := ToTime(testtime)

	// assert
	assert.Equal(t, expectedtime, gotime, "StringTime error")
}

func TestParseTMPeriod(t *testing.T) {
	// todo!
}

func TestToPostgresPeriodFormat(t *testing.T) {
	// todo!
}
