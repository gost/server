package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"encoding/json"
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
	assert.True(t,len(propertynames)>0)
	assert.True(t,len(supportedencoding)==0)
}

func TestSetLinksAdvanced(t *testing.T) {
	// arrange
	thing := &Thing{}
	sensor := &Sensor{}
	datastream := &Datastream{}
	observedproperty := &ObservedProperty{}
	// todo: how to set observations???
	// obs1 := &Observation{}
	// obs2 := &Observation{}
	//obs := []Observation{}
	//datastream.Observations=*obs

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



func TestParseEntity(t *testing.T){
	// arrange
	datastream := &Datastream{}
	datastream.ID = "0"
	dsjson,_:=json.Marshal(datastream)

	// act
	datastream.ParseEntity(dsjson)

	// assert
	assert.True(t,datastream.ID == "0")
}

func TestGetEntityType(t *testing.T) {
	// arrange
	datastream := &Datastream{}

	// act
	entityType := datastream.GetEntityType()

	// assert
	assert.Equal(t, EntityTypeDatastream, entityType, "GetEntityType should be Datastream")
}

func TestContainsMandatoryParameters(t *testing.T) {
	// arrange
	datastream := &Datastream{}

	// act
	contains, _ := datastream.ContainsMandatoryParams()

	// assert
	assert.False(t, contains, "Datastream is expected not to have mandatory paramaters")
}
