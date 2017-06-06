package rest

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEndPointDatastream(t *testing.T) {
	// arrange
	ep := createDatastreamsEndpoint("http://www.nu.nl")
	ep.Name = "yo"

	// assert
	assert.True(t, ep != nil)
	assert.True(t, ep.GetName() == "yo")
}
