package entities

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
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
	assert.Equal(t, "ObservedProperties", EntityLinkObservedPropertys.ToString())
	assert.Equal(t, "Observations", EntityLinkObservations.ToString())
	assert.Equal(t, "FeatureOfInterests", EntityLinkFeatureOfInterests.ToString())
}

func TestCreateEntitySelfLink(t *testing.T) {
	//act
	selfLink := CreateEntitySelfLink(externalURL, lt.ToString(), "")
	selfLinkWithID := CreateEntitySelfLink(externalURL, lt.ToString(), id)

	//assert
	assert.Equal(t, fmt.Sprintf("%s/v1.0/Things", externalURL), selfLink, "Entityselflink is not in the correct format")
	assert.Equal(t, fmt.Sprintf("%s/v1.0/Things(myid)", externalURL), selfLinkWithID, "Entityselflink with id is not in the correct format")
}

func TestCreateEntityLink(t *testing.T) {
	//act
	link := CreateEntityLink(true, externalURL, lt.ToString(), ls.ToString(), "")
	linkWithID := CreateEntityLink(true, externalURL, lt.ToString(), ls.ToString(), id)
	linkEmpty := CreateEntityLink(false, externalURL, lt.ToString(), ls.ToString(), "")

	//assert
	assert.Equal(t, fmt.Sprintf("%s/v1.0/%s/%s", externalURL, lt.ToString(), ls.ToString()), link, "EntityLink is not in the correct format")
	assert.Equal(t, fmt.Sprintf("%s/v1.0/%s(%s)/%s", externalURL, lt.ToString(), id, ls.ToString()), linkWithID, "EntityLink with id is not in the correct format")
	assert.Equal(t, "", linkEmpty, "EntityLink link should be empty")
}

func TestCheckMandatoryParamNoErrors(t *testing.T) {
	//arrange
	errLis1 := []error{}
	errLis2 := []error{}
	errLis3 := []error{}
	errLis4 := []error{}
	errLis5 := []error{}

	testString := "test"
	testMap := map[string]string{"test": "test"}

	testThing := &Sensor{}
	testThing.ID = "1"

	testSensor := &Sensor{}
	testSensor.ID = "1"

	testObservedProperty := &ObservedProperty{}
	testObservedProperty.ID = "1"

	//act
	CheckMandatoryParam(&errLis1, testString, et, "test")
	CheckMandatoryParam(&errLis2, testMap, et, "test")
	CheckMandatoryParam(&errLis3, testThing, et, "test")
	CheckMandatoryParam(&errLis4, testSensor, et, "test")
	CheckMandatoryParam(&errLis5, testObservedProperty, et, "test")

	//assert
	assert.Equal(t, len(errLis1), 0, "CheckMandatoryParam string should not have returned an error")
	assert.Equal(t, len(errLis2), 0, "CheckMandatoryParam map[string]string should not have returned an error")
	assert.Equal(t, len(errLis3), 0, "CheckMandatoryParam Thing should not have returned an error")
	assert.Equal(t, len(errLis4), 0, "CheckMandatoryParam Sensor should not have returned an error")
	assert.Equal(t, len(errLis5), 0, "CheckMandatoryParam ObservedProperty should not have returned an error")
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

func TestEncodingToString(t *testing.T) {
	//assert
	assert.Equal(t, "unknown", EncodingUnknown.Value)
	assert.Equal(t, 0, EncodingUnknown.Code, "EncodingUnknown code changed")

	assert.Equal(t, "application/vnd.geo+json", EncodingGeoJSON.Value)
	assert.Equal(t, 1, EncodingGeoJSON.Code, "EncodingGeoJSON code changed")

	assert.Equal(t, "application/pdf", EncodingPDF.Value)
	assert.Equal(t, 2, EncodingPDF.Code, "EncodingPDF code changed")

	assert.Equal(t, "http://www.opengis.net/doc/IS/SensorML/2.0", EncodingSensorML.Value)
	assert.Equal(t, 3, EncodingSensorML.Code, "EncodingSensorML code changed")
}

func TestEncodingTypeOk(t *testing.T) {
	//arrange
	sml := "http://www.opengis.net/doc/IS/SensorML/2.0"

	//act
	encoding, err := CreateEncodingType(sml)

	//assert
	assert.Nil(t, err, fmt.Sprintf("Creating encoding type for %s should not have returned an error", sml))
	assert.Equal(t, 3, encoding.Code, fmt.Sprintf("Incorrect encoding code for %s", sml))
}

func TestEncodingTypeFail(t *testing.T) {
	//arrange
	sml := "http://www.opengis.net/doc/IS/SensorM/2.0"

	//act
	_, err := CreateEncodingType(sml)

	//assert
	assert.NotNil(t, err, fmt.Sprintf("Creating encoding type for %s should not returned an error", sml))
}

func TestCheckEncodingSupportedSensorOk(t *testing.T) {
	//arrange
	s := &Sensor{}

	//act
	_, err := CheckEncodingSupported(s, "http://www.opengis.net/doc/IS/SensorML/2.0")

	//assert
	assert.Nil(t, err, "Sensor should support encoding http://www.opengis.net/doc/IS/SensorML/2.0")
}

func TestCheckEncodingSupportedSensorFail(t *testing.T) {
	//arrange
	s := &Sensor{}

	//act
	_, err := CheckEncodingSupported(s, "http://www.opengis.net/doc/IS/SensorML/2")

	//assert
	assert.NotNil(t, err, "Sensor should not support encoding http://www.opengis.net/doc/IS/SensorML/2")
}
