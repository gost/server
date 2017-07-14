package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEndPointSensor(t *testing.T) {
	// arrange
	ep := CreateSensorsEndpoint("http://www.nu.nl")
	ep.Name = "yo"

	// assert
	assert.True(t, ep != nil)
	assert.True(t, ep.GetName() == "yo")
}
