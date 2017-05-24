package entities

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	//act
	_, err := CheckEncodingSupported("http://www.opengis.net/doc/IS/SensorML/2.0")

	//assert
	assert.Nil(t, err, "Sensor should support encoding http://www.opengis.net/doc/IS/SensorML/2.0")
}

func TestCheckEncodingSupportedSensorFail(t *testing.T) {
	//arrange

	//act
	_, err := CheckEncodingSupported("http://www.opengis.net/doc/IS/SensorML/2")

	//assert
	assert.NotNil(t, err, "Sensor should not support encoding http://www.opengis.net/doc/IS/SensorML/2")
}
