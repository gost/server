package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHistoricalLocationWithMandatoryParameters(t *testing.T) {
	//arrange
	historicalLocation := &HistoricalLocation{}
	historicalLocation.Time = "testtime"
	historicalLocation.Locations = []*Location{}
	thing := Thing{}
	thing.ID = 1
	historicalLocation.Thing = &thing

	location := &Location{}
	location.ID = 1
	historicalLocation.Locations = append(historicalLocation.Locations, location)
	//act
	res, err := historicalLocation.ContainsMandatoryParams()

	//assert
	assert.True(t, res, "HistoricalLocation result should be true")
	assert.Nil(t, err, "HistoricalLocation mandatory param description not filled in should have returned error")
}

func TestHistoricalLocationWithoutMandatoryParameters(t *testing.T) {
	//arrange
	historicalLocation := &HistoricalLocation{}

	//act
	res, err := historicalLocation.ContainsMandatoryParams()

	//assert
	assert.False(t, res, "HistoricalLocation result should be false")
	assert.NotNil(t, err, "HistoricalLocation mandatory param description not filled in should have returned error")
}

func TestSetAllLinks(t *testing.T) {
	//arrange
	historicalLocation := &HistoricalLocation{}
	historicalLocation.Time = "testtime"
	thing := Thing{}
	thing.Description = "testdescription"
	historicalLocation.Thing = &thing

	//act
	historicalLocation.SetAllLinks("http://www.test.com")

	//assert
	assert.NotNil(t, historicalLocation.NavSelf)
	assert.NotNil(t, historicalLocation.NavThing)
}

func TestParseHistoricalLocationShouldFail(t *testing.T) {
	//arrange
	historicalLocation := &HistoricalLocation{}

	//act
	err := historicalLocation.ParseEntity([]byte("hallo"))

	//assert
	assert.NotEqual(t, err, nil, "Historical parse from json should have failed")
}
