package entities

import (
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

func TestCreateEntitySelfLink(t *testing.T) {
	//act
	selfLink := CreateEntitySelfLink(externalURL, lt.ToString(), "")
	selfLinkWithID := CreateEntitySelfLink(externalURL, lt.ToString(), id)

	//assert
	assert.Equal(t, "www.myurl.nl/v1.0/Things", selfLink, "Entityselflink is not in the correct format")
	assert.Equal(t, "www.myurl.nl/v1.0/Things(myid)", selfLinkWithID, "Entityselflink with id is not in the correct format")
}

func TestCreateEntityLink(t *testing.T) {
	//act
	selfLink := CreateEntityLink(true, lt.ToString(), ls.ToString(), "")
	selfLinkWithID := CreateEntityLink(true, lt.ToString(), ls.ToString(), id)

	//assert
	assert.Equal(t, "../Things/Sensors", selfLink, "EntityLink is not in the correct format")
	assert.Equal(t, "../Things(myid)/Sensors", selfLinkWithID, "EntityLink with id is not in the correct format")
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
	testThing := &Thing{ID: "1"}
	testSensor := &Sensor{ID: "1"}
	testObservedProperty := &ObservedProperty{ID: "1"}

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
	testThingEmpty := &Thing{}
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
