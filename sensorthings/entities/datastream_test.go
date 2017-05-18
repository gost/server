package entities

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetLinksReturnsVariousLinks(t *testing.T) {
	// arrange
	datastream := &Datastream{}
	datastream.ID = "0"
	// todo: what happens if there is no ID?

	// act
	datastream.SetAllLinks("http://www.nu.nl")
	propertynames := datastream.GetPropertyNames()
	supportedencoding := datastream.GetSupportedEncoding()

	// assert
	assert.NotNil(t, datastream.NavSelf, "NavSelf should be filled in")
	assert.NotNil(t, datastream.NavThing, "NavThing should be filled in")
	assert.NotNil(t, datastream.NavSensor, "NavSensor should be filled in")
	assert.NotNil(t, datastream.NavObservations, "NavObservations should be filled in")
	assert.NotNil(t, datastream.NavObservedProperty, "NavObservedProperty should be filled in")
	assert.True(t, len(propertynames) > 0)
	assert.True(t, len(supportedencoding) == 0)
}

func TestSetLinksAdvanced(t *testing.T) {
	// arrange
	thing := &Thing{}
	sensor := &Sensor{}
	datastream := &Datastream{}
	observedproperty := &ObservedProperty{}
	obs1 := &Observation{}
	obs2 := &Observation{}
	obs := []*Observation{obs1, obs2}
	datastream.Observations = obs

	datastream.ID = "0"
	datastream.Thing = thing
	datastream.Sensor = sensor
	datastream.ObservedProperty = observedproperty

	// act
	datastream.SetAllLinks("http://www.nu.nl")

	// assert
	assert.NotNil(t, datastream.NavSelf, "NavSelf should be filled in")
	assert.NotNil(t, datastream.NavThing, "NavThing should be filled in")
	assert.NotNil(t, datastream.NavSensor, "NavSensor should be filled in")
	assert.NotNil(t, datastream.NavObservations, "NavObservations should be filled in")
	assert.NotNil(t, datastream.NavObservedProperty, "NavObservedProperty should be filled in")
}

func TestParseEntity(t *testing.T) {
	// arrange
	datastream := &Datastream{}
	datastream.ID = "0"
	dsjson, _ := json.Marshal(datastream)

	// act
	datastream.ParseEntity(dsjson)

	// assert
	assert.True(t, datastream.ID == "0")
}

func TestParseEntityFails(t *testing.T) {
	// arrange
	dsjson, _ := json.Marshal("hoho")

	// act
	datastream := &Datastream{}
	err := datastream.ParseEntity(dsjson)

	// assert
	assert.True(t, err != nil)
}

func TestGetEntityType(t *testing.T) {
	// arrange
	datastream := &Datastream{}

	// act
	entityType := datastream.GetEntityType()

	// assert
	assert.Equal(t, EntityTypeDatastream, entityType, "GetEntityType should be Datastream")
}

func TestContainsMandatoryParametersFails(t *testing.T) {
	// arrange
	datastream := &Datastream{}

	// act
	contains, _ := datastream.ContainsMandatoryParams()

	// assert
	assert.False(t, contains, "Datastream is expected not to have mandatory paramaters")
}

func TestContainsMandatoryParametersSucceeds(t *testing.T) {
	// arrange
	datastream := &Datastream{}
	datastream.Name = "name"
	datastream.Description = "desc"
	unitofmeasurement := map[string]interface{}{"definition": "http://www.qudt.org/qudt/owl/1.0.0/unit/Instances.html/Lumen", "name": "Centigrade", "symbol": "C"}
	datastream.UnitOfMeasurement = unitofmeasurement
	datastream.ObservationType = "OMCategoryUnknown"
	thing := &Thing{}
	thing.Name = "name"
	thing.Description = "desc"
	datastream.Thing = thing
	sensor := &Sensor{
		Name:         "test",
		Description:  "test",
		EncodingType: "test",
		Metadata:     "test",
	}
	datastream.Sensor = sensor
	op := &ObservedProperty{
		Name:        "test",
		Definition:  "test",
		Description: "test",
	}


	datastream.ObservedProperty = op

	// act
	contains, _ := datastream.ContainsMandatoryParams()

	// assert
	assert.True(t, contains, "Datastream is expected to have mandatory paramaters")
}
