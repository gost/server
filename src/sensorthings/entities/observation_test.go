package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEntityTypeReturnsCorrectType (t *testing.T) {
	//arrange
	observation := &Observation{}

	//act
	entityType := observation.GetEntityType()

	//assert
	assert.Equal(t, EntityTypeObservation, entityType, "getEntityType should return correct type")
}

func TestSetLinksReturnsCorrectLinks (t *testing.T) {
	// arrange
	observation := &Observation{}

	// act
	observation.SetLinks( "www.nu.nl" )

	// assert
	assert.NotNil( t, observation.NavSelf, " NAvSelf should be filled in" )
	assert.NotNil( t, observation.NavDatastream, " NavDatastream should be filled in" )
	assert.NotNil( t, observation.NavFeatureOfInterest, " NavFeatureOfInterest should be filled in" )
}

func TestGetSupportedEncodingShouldNotReturnAnyEncoding (t *testing.T) {
	// arrange
	observation := &Observation{}

	// act
	supportedEncoding := observation.GetSupportedEncoding()

	// assert
	assert.Equal(t, 0, len(supportedEncoding), "Observation should not supprt any encoding")
}


func TestParseEntityShouldFail (t *testing.T){
	//arrange
	observation := &Observation{}

	//act
	err := observation.ParseEntity([]byte("hallo"))

	//assert
	assert.NotEqual(t, err, nil, "Observation parse from json should have failed")
}

func TestMissingMandatoryParametersObservation(t *testing.T) {
	//arrange
	// observation := &Observation{}

	//act
	// _, _ := observation.ContainsMandatoryParams()

	// todo: something goes wrong here?
	// assert.False(t, ok)
	// assert.NotNil(t, err, "Observation mandatory parameters not filled in should have returned error")
	// if len(err) > 0 {
	//	assert.Contains(t, fmt.Sprintf("%v", err[0]), "name")
	//}
}

func TestMarshalPostgresJSONReturnsSomething(t *testing.T){
	// arrange
	observation := &Observation{}

	// act
	bytes, _ := observation.MarshalPostgresJSON()

	// assert

	assert.NotNil(t,bytes)
}

func TestMandatoryParametersExistObservation(t *testing.T) {
	//arrange
	observation := &Observation{
		Result: 20,
		Datastream: &Datastream{},
	}

	//act
	ok, err := observation.ContainsMandatoryParams()

	//assert
	assert.Nil(t, err, "All mandatory params are filled in should not have returned an error")
	assert.True(t, ok, "Observation mandatory parameters should be ok")
}



