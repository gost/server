package rest

import (
"testing"
"github.com/stretchr/testify/assert"
)

func TestGetEndPointObservedProperties(t *testing.T) {
	// arrange
	ep := createObservedPropertiesEndpoint("http://www.nu.nl")
	ep.Name = "yo"

	// assert
	assert.True(t, ep != nil)
	assert.True(t, ep.GetName() == "yo")
}
