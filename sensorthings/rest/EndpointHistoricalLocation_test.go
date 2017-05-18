package rest

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetEndPointHistoricalLocation(t *testing.T) {
	// arrange
	ep := createHistoricalLocationsEndpoint("http://www.nu.nl")
	ep.Name = "yo"

	// assert
	assert.True(t, ep != nil)
	assert.True(t, ep.GetName() == "yo")
}

