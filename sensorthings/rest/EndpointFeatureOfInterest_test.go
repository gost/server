package rest

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetEndPointFoi(t *testing.T) {
	// arrange
	ep := createFeaturesOfInterestEndpoint("http://www.nu.nl")
	ep.Name = "yo"

	// assert
	assert.True(t, ep != nil)
	assert.True(t, ep.GetName() == "yo")
}
