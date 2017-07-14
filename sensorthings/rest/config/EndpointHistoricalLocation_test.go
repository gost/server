package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEndPointHistoricalLocation(t *testing.T) {
	// arrange
	ep := CreateHistoricalLocationsEndpoint("http://www.nu.nl")
	ep.Name = "yo"

	// assert
	assert.True(t, ep != nil)
	assert.True(t, ep.GetName() == "yo")
}
