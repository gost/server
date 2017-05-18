package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestObservationTypesShouldReturnArray(t *testing.T) {
	// assert
	assert.True(t, len(ObservationTypes) > 0, "Array of Observationtypes should be returned")
}

func TestGetObservationTypeByValueShouldReturnCorrectObservationType(t *testing.T) {
	// arrange
	input := "http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_CategoryObservation"

	// act
	res, err := GetObservationTypeByValue(input)

	// assert
	assert.Nil(t, err, "Error should not be nil")
	assert.True(t, res.Code == 1, "Code should be 1")
}

func TestGetObservationTypeByCodeShouldReturnCorrectObservationType(t *testing.T) {
	// arrange
	expected := "http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_CategoryObservation"
	input := int64(1)

	// act
	res, err := GetObservationTypeByID(input)

	// assert
	assert.Nil(t, err, "Error should be nil")
	assert.True(t, res.Value == expected, "Code should be "+expected)
}

func TestGetObservationTypeByNotExistingCodeShouldReturnError(t *testing.T) {
	// arrange
	expected := "http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_CategoryObservation"
	input := int64(9999)

	// act
	res, err := GetObservationTypeByID(input)

	// assert
	assert.NotNil(t, err, "Error should not be nil")
	assert.True(t, res.Value == expected, "Code should be "+expected)
}

func TestGetObservationTypeByNotExistingValueShouldReturnError(t *testing.T) {
	// arrange
	input := "not_existing"
	expected := "http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_CategoryObservation"

	// act
	res, err := GetObservationTypeByValue(input)

	// assert
	assert.NotNil(t, err, "Error should not be nil")
	assert.True(t, res.Value == expected, "Code should be "+expected)
}
