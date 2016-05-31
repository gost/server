package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestObservationTypessShouldReturnArray(t *testing.T) {
	// assert
	assert.True(t,len(ObservationTypes)>0, "Array of Observationtypes should be returned")
}
