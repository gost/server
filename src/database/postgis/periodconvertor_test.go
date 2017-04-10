package postgis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//TestIso8601 tests the ToIso8601 function
func TestIso8601(t *testing.T) {
	// arrange
	var testdate = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	// act
	isodate := ToIso8601(testdate)
	// assert
	assert.Equal(t, "2009-11-10T23:00:00.000Z", isodate, "Iso8601 error")
}

//TestToTime tests the ToTime function
func TestToTime(t *testing.T) {
	// arrange
	var testtime = "2014-03-01T13:00:00Z"
	var expectedtime = time.Date(2014, time.March, 1, 13, 0, 0, 0, time.UTC)
	// act
	gotime, _ := ToTime(testtime)

	// assert
	assert.Equal(t, expectedtime, gotime, "StringTime error")
}

// TestParseTMPeriod tests the ParseTMPeriod function
func TestParseTMPeriod(t *testing.T) {
	// arrange
	var period = "2014-03-01T13:00:00Z/2015-05-11T15:30:00Z"
	var expectedStart = time.Date(2014, time.March, 1, 13, 0, 0, 0, time.UTC)
	var expectedEnd = time.Date(2015, time.May, 11, 15, 30, 0, 0, time.UTC)

	// act
	var actualPeriod = ParseTMPeriod(period)

	// assert
	assert.Equal(t, expectedStart, actualPeriod[0], "Parse period start error")
	assert.Equal(t, expectedEnd, actualPeriod[1], "Parse period end error")

}

// TestToPostgresPeriodFormat tests the ToPostgresPeriodFormat function
func TestToPostgresPeriodFormat(t *testing.T) {
	// arrange
	var period [2]time.Time
	period[0] = time.Date(2014, time.March, 1, 13, 0, 0, 0, time.UTC)
	period[1] = time.Date(2015, time.May, 11, 15, 30, 0, 0, time.UTC) 
	var expectedResult = "[2014-03-01T13:00:00.000Z,2015-05-11T15:30:00.000Z]"

	// act
	var postgresPeriod = ToPostgresPeriodFormat(period)

	// assert
	assert.Equal(t,expectedResult,postgresPeriod,"ToPostgresPeriodFormat error")
}
