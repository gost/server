package entities

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var jsonSensor = `{
    "name": "MySensor",
    "description": "SensorUp Tempomatic 2000",
    "encodingType": "http://schema.org/description",
    "metadata": "Calibration date:  Jan 1, 2014"
}`

var jsonSensorError = `{
    "description": "SensorUp Tempomatic 2000",
}`

func TestMissingMandatoryParametersSensor(t *testing.T) {
	//arrange
	sensor := &Sensor{}

	//act
	_, err := sensor.ContainsMandatoryParams()

	//assert
	assert.NotNil(t, err, "Sensor mandatory param description not filled in should have returned error")
	if len(err) > 0 {
		assert.Contains(t, fmt.Sprintf("%v", err[0]), "name")
	}
}

func TestSensorGetPropertyNames(t *testing.T) {
	// arrange
	sensor := &Sensor{}

	// act
	propertynames := sensor.GetPropertyNames()

	// assert
	assert.True(t, propertynames[0] == "id")
}

func TestMandatoryParametersExistSensor(t *testing.T) {
	//arrange
	sensor := &Sensor{
		Name:         "test",
		Description:  "test",
		EncodingType: "test",
		Metadata:     "test",
	}

	//act
	_, err := sensor.ContainsMandatoryParams()

	//assert
	assert.Nil(t, err, "All mandatory params are filled in should not have returned an error")
}

func TestParseEntityResultOkSensor(t *testing.T) {
	//arrange
	sensor := &Sensor{}

	//act
	err := sensor.ParseEntity([]byte(jsonSensor))

	//assert
	assert.Equal(t, err, nil, "Unable to parse json into thing")
}

func TestParseEntityResultNotOkSensor(t *testing.T) {
	//arrange
	thing := &Sensor{}

	//act
	err := thing.ParseEntity([]byte(jsonSensorError))

	//assert
	assert.NotEqual(t, err, nil, "Sensor parse from json should have failed")
}

func TestSetLinksSensor(t *testing.T) {
	//arrange
	sensor := &Sensor{}
	sensor.ID = id
	ds1 := &Datastream{}
	ds2 := &Datastream{}
	datastreams := []*Datastream{ds1, ds2}
	sensor.Datastreams = datastreams

	//act
	sensor.SetAllLinks(externalURL)

	//assert
	assert.Equal(t, sensor.NavSelf, fmt.Sprintf("%s/v1.0/%s(%s)", externalURL, EntityLinkSensors.ToString(), id), "Sensor navself incorrect")
	//assert.Equal(t, sensor.NavDatastreams, fmt.Sprintf("%s/v1.0/%s(%s)/%s", externalURL, EntityLinkSensors.ToString(), id, EntityLinkDatastreams.ToString()), "Sensor NavDatastreams incorrect")
}
