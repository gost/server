package entities

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	lt          = EntityLinkThings
	ls          = EntityLinkSensors
	et          = EntityTypeThing
	externalURL = "www.myurl.nl"
	id          = "myid"
)

func TestEntityTypeStrings(t *testing.T) {
	assert.Equal(t, "Thing", EntityTypeThing.ToString())
	assert.Equal(t, "Location", EntityTypeLocation.ToString())
	assert.Equal(t, "HistoricalLocation", EntityTypeHistoricalLocation.ToString())
	assert.Equal(t, "Datastream", EntityTypeDatastream.ToString())
	assert.Equal(t, "Sensor", EntityTypeSensor.ToString())
	assert.Equal(t, "ObservedProperty", EntityTypeObservedProperty.ToString())
	assert.Equal(t, "Observation", EntityTypeObservation.ToString())
	assert.Equal(t, "FeatureOfInterest", EntityTypeFeatureOfInterest.ToString())
}

func TestEntityLinkStrings(t *testing.T) {
	assert.Equal(t, "Things", EntityLinkThings.ToString())
	assert.Equal(t, "Locations", EntityLinkLocations.ToString())
	assert.Equal(t, "HistoricalLocations", EntityLinkHistoricalLocations.ToString())
	assert.Equal(t, "Datastreams", EntityLinkDatastreams.ToString())
	assert.Equal(t, "Sensors", EntityLinkSensors.ToString())
	assert.Equal(t, "ObservedProperties", EntityLinkObservedProperties.ToString())
	assert.Equal(t, "Observations", EntityLinkObservations.ToString())
	assert.Equal(t, "FeatureOfInterest", EntityLinkFeatureOfInterests.ToString())
}

func TestCreateEntitySelfLink(t *testing.T) {
	//act
	selfLink := CreateEntitySelfLink(externalURL, lt.ToString(), nil)
	selfLinkWithID := CreateEntitySelfLink(externalURL, lt.ToString(), id)

	//assert
	assert.Equal(t, fmt.Sprintf("%s/v1.0/Things", externalURL), selfLink, "Entityselflink is not in the correct format")
	assert.Equal(t, fmt.Sprintf("%s/v1.0/Things(myid)", externalURL), selfLinkWithID, "Entityselflink with id is not in the correct format")
}

func TestCreateEntityLink(t *testing.T) {
	//act
	link := CreateEntityLink(true, externalURL, lt.ToString(), ls.ToString(), nil)
	linkWithID := CreateEntityLink(true, externalURL, lt.ToString(), ls.ToString(), id)
	linkEmpty := CreateEntityLink(false, externalURL, lt.ToString(), ls.ToString(), nil)

	//assert
	assert.Equal(t, fmt.Sprintf("%s/v1.0/%s/%s", externalURL, lt.ToString(), ls.ToString()), link, "EntityLink is not in the correct format")
	assert.Equal(t, fmt.Sprintf("%s/v1.0/%s(%s)/%s", externalURL, lt.ToString(), id, ls.ToString()), linkWithID, "EntityLink with id is not in the correct format")
	assert.Equal(t, "", linkEmpty, "EntityLink link should be empty")
}

func TestBaseEntity(t *testing.T) {
	// arrange
	b := BaseEntity{}

	// act
	parse := b.ParseEntity(nil)
	m, _ := b.ContainsMandatoryParams()
	linkerror := b.SetLinks("")
	id := b.GetID()
	entitytype := b.GetEntityType()
	propnames := b.GetPropertyNames()
	selflink := b.GetSelfLink()
	encoding := b.GetSupportedEncoding()
	b.SetID(35)

	// assert
	assert.Nil(t, parse)
	assert.False(t, m)
	assert.Nil(t, linkerror)
	assert.Nil(t, id)
	assert.Nil(t, propnames)
	assert.True(t, entitytype == EntityTypeUnknown)
	assert.True(t, selflink == "")
	assert.Nil(t, encoding)

}

func TestCheckMandatoryParamNoErrors(t *testing.T) {
	//arrange
	errLis1 := []error{}
	errLis2 := []error{}
	errLis3 := []error{}
	errLis4 := []error{}
	errLis5 := []error{}
	errLis6 := []error{}
	errLis7 := []error{}
	errLis8 := []error{}
	errLis9 := []error{}
	errLis10 := []error{}
	errLis11 := []error{}
	errLis12 := []error{}

	testString := "test"
	testMap := map[string]string{"test": "test"}
	emptytestMap := map[string]string{}

	testemptystring := ""

	testThing := &Thing{}
	testThing.ID = "1"

	testSensor := &Sensor{}
	testSensor.ID = "1"

	testObservedProperty := &ObservedProperty{}
	testObservedProperty.ID = "1"

	datastream := &Datastream{}
	datastream.ID = "3"

	var nilstring *string
	var datastreamnil *Datastream

	//act
	CheckMandatoryParam(&errLis1, testString, et, "test")
	CheckMandatoryParam(&errLis2, testMap, et, "test")
	CheckMandatoryParam(&errLis3, testThing, et, "test")
	CheckMandatoryParam(&errLis4, testSensor, et, "test")
	CheckMandatoryParam(&errLis5, testObservedProperty, et, "test")
	CheckMandatoryParam(&errLis6, &testString, et, "test")
	CheckMandatoryParam(&errLis7, datastream, et, "test")
	CheckMandatoryParam(&errLis8, nil, et, "test")
	CheckMandatoryParam(&errLis9, &testemptystring, et, "test")
	CheckMandatoryParam(&errLis10, nilstring, et, "test")
	CheckMandatoryParam(&errLis11, datastreamnil, et, "test")
	CheckMandatoryParam(&errLis12, emptytestMap, et, "test")

	//assert
	assert.Equal(t, len(errLis1), 0, "CheckMandatoryParam string should not have returned an error")
	assert.Equal(t, len(errLis2), 0, "CheckMandatoryParam map[string]string should not have returned an error")
	assert.Equal(t, len(errLis3), 0, "CheckMandatoryParam Thing should not have returned an error")
	assert.Equal(t, len(errLis4), 0, "CheckMandatoryParam Sensor should not have returned an error")
	assert.Equal(t, len(errLis5), 0, "CheckMandatoryParam ObservedProperty should not have returned an error")
	assert.Equal(t, len(errLis6), 0, "CheckMandatoryParam string pointer should not have returned an error")
	assert.Equal(t, len(errLis7), 0, "CheckMandatoryParam datastream pointer should not have returned an error")
	assert.Equal(t, len(errLis8), 1, "CheckMandatoryParam nil parameter should have returned an error")
	assert.Equal(t, len(errLis9), 1, "CheckMandatoryParam empty string  parameter should have returned an error")
	assert.Equal(t, len(errLis10), 1, "CheckMandatoryParam nil string pointer parameter should have returned an error")
	assert.Equal(t, len(errLis11), 1, "CheckMandatoryParam datastream nil pointer parameter should have returned an error")
	assert.Equal(t, len(errLis12), 1, "CheckMandatoryParam empty string map parameter should have returned an error")
}

func TestCheckMandatoryParamErrors(t *testing.T) {
	//arrange
	errLis1 := []error{}
	errLis2 := []error{}
	errLis3 := []error{}
	errLis4 := []error{}
	errLis5 := []error{}

	testStringEmpty := ""
	testMapEmpty := map[string]string{}
	testThingEmpty := &Sensor{}
	testSensorEmpty := &Sensor{}
	testObservedPropertyEmpty := &ObservedProperty{}

	//act
	CheckMandatoryParam(&errLis1, testStringEmpty, et, "test")
	CheckMandatoryParam(&errLis2, testMapEmpty, et, "test")
	CheckMandatoryParam(&errLis3, testThingEmpty, et, "test")
	CheckMandatoryParam(&errLis4, testSensorEmpty, et, "test")
	CheckMandatoryParam(&errLis5, testObservedPropertyEmpty, et, "test")

	//assert
	assert.Len(t, errLis1, 1, "CheckMandatoryParam string should have returned an error")
	assert.Len(t, errLis2, 1, "CheckMandatoryParam map[string]string should have returned an error")
	assert.Len(t, errLis3, 1, "CheckMandatoryParam Thing should have returned an error")
	assert.Len(t, errLis4, 1, "CheckMandatoryParam Sensor should have returned an error")
	assert.Len(t, errLis5, 1, "CheckMandatoryParam ObservedProperty should have returned an error")
}
