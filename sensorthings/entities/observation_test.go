package entities

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var jsonObservation = ` {
         "@iot.id": 1368606,
         "@iot.selfLink": "http://gost.geodan.nl/v1.0/Observations(1368606)",
         "phenomenonTime": "2017-05-08T09:09:52.032311769Z",
         "result": 10.291769996296797,
         "Datastream@iot.navigationLink": "http://gost.geodan.nl/v1.0/Observations(1368606)/Datastream",
         "FeatureOfInterest@iot.navigationLink": "http://gost.geodan.nl/v1.0/Observations(1368606)/FeatureOfInterest",
         "resultTime": null
      }`

func TestGetEntityTypeReturnsCorrectType(t *testing.T) {
	//arrange
	observation := &Observation{}

	//act
	entityType := observation.GetEntityType()

	//assert
	assert.Equal(t, EntityTypeObservation, entityType, "getEntityType should return correct type")
}

func TestSetLinksReturnsCorrectLinks(t *testing.T) {
	// arrange
	observation := &Observation{}
	datastream := &Datastream{}
	foi := &FeatureOfInterest{}

	observation.Datastream = datastream
	observation.FeatureOfInterest = foi

	// act
	observation.SetAllLinks("www.nu.nl")

	// assert
	assert.NotNil(t, observation.NavSelf, " NAvSelf should be filled in")
	assert.NotNil(t, observation.NavDatastream, " NavDatastream should be filled in")
	assert.NotNil(t, observation.NavFeatureOfInterest, " NavFeatureOfInterest should be filled in")
	assert.NotNil(t, observation.FeatureOfInterest.NavSelf, " NavSelf FeatureofInterest should be filled in")
	assert.NotNil(t, observation.Datastream.NavSelf, " NavSelf Datastream should be filled in")

}

func TestMarshalJson(t *testing.T) {
	// arrange
	observation := &Observation{}

	// act
	obsjson, _ := observation.MarshalJSON()

	// assert
	assert.NotNil(t, obsjson)
}

func TestMarshalJsonWithTimes(t *testing.T) {
	// arrange
	observation := &Observation{}
	resultTime := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	observation.ResultTime = &resultTime

	// act
	obsjson, _ := observation.MarshalJSON()

	// assert
	assert.NotNil(t, obsjson)
}

func TestMarshalJsonWithEmptyTimes(t *testing.T) {
	// arrange
	observation := &Observation{}
	resultTime := ""
	observation.ResultTime = &resultTime

	// act
	obsjson, _ := observation.MarshalJSON()

	// assert
	assert.NotNil(t, obsjson)
}

func TestGetPropertyNames(t *testing.T) {
	// arrange
	observation := &Observation{}

	// act
	propertynames := observation.GetPropertyNames()

	// assert
	assert.True(t, propertynames[0] == "id")
}

func TestGetSupportedEncodingShouldNotReturnAnyEncoding(t *testing.T) {
	// arrange
	observation := &Observation{}

	// act
	supportedEncoding := observation.GetSupportedEncoding()

	// assert
	assert.Equal(t, 0, len(supportedEncoding), "Observation should not supprt any encoding")
}

func TestParseEntityShouldWork(t *testing.T) {
	//arrange
	observation := &Observation{}

	//act
	err := observation.ParseEntity([]byte(jsonObservation))

	//assert
	assert.Equal(t, err, nil, "Observation parse from json should have succeeded")
}

func TestParseEntityShouldFail(t *testing.T) {
	//arrange
	observation := &Observation{}

	//act
	err := observation.ParseEntity([]byte("hallo"))

	//assert
	assert.NotEqual(t, err, nil, "Observation parse from json should have failed")
}

func TestMissingMandatoryParametersObservation(t *testing.T) {
	//arrange
	observation := &Observation{}

	//act
	ok, err := observation.ContainsMandatoryParams()

	assert.False(t, ok)
	assert.NotNil(t, err, "Observation mandatory parameters not filled in should have returned error")
	if len(err) > 0 {
		assert.Contains(t, fmt.Sprintf("%v", err[0]), "result")
	}
}

func TestMissingMandatoryParametersWithTimesObservation(t *testing.T) {
	//arrange
	observation := &Observation{}
	observation.PhenomenonTime = time.Now().UTC().Format(time.RFC3339Nano)
	resultTime := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	observation.ResultTime = &resultTime

	//act
	ok, _ := observation.ContainsMandatoryParams()

	// assert
	assert.False(t, ok)
}

func TestMissingMandatoryParametersWithWrongTimesObservation(t *testing.T) {
	//arrange
	observation := &Observation{}
	observation.PhenomenonTime = "haha"
	resulttime := "hoho"
	observation.ResultTime = &resulttime

	//act
	ok, _ := observation.ContainsMandatoryParams()

	// assert
	assert.False(t, ok)
}

func TestMarshalPostgresJSONReturnsSomething(t *testing.T) {
	// arrange
	observation := &Observation{}

	// act
	bytes, _ := observation.MarshalPostgresJSON()

	// assert

	assert.NotNil(t, bytes)
}

func TestMarshalPostgresJSONWithResultTime(t *testing.T) {
	// arrange
	observation := &Observation{}
	resultTime := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	observation.ResultTime = &resultTime

	// act
	bytes, _ := observation.MarshalPostgresJSON()

	// assert

	assert.NotNil(t, bytes)
}

func TestMandatoryParametersExistObservation(t *testing.T) {
	//arrange
	observation := &Observation{
		Result:     20,
		Datastream: &Datastream{},
	}
	observation.Datastream.ID = "1"

	//act
	ok, err := observation.ContainsMandatoryParams()

	//assert
	assert.Nil(t, err, "All mandatory params are filled in should not have returned an error")
	assert.True(t, ok, "Observation mandatory parameters should be ok")
}
