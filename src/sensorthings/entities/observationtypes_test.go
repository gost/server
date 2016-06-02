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
	assert.True(t, res.Code == 0, "Code should be 0")
}
