package rest

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetEndPointVersion(t *testing.T) {
	// arrange
	ep := createVersionEndpoint("http://www.nu.nl")
	ep.Name = "yo"

	// assert
	assert.True(t, ep.GetName() == "yo")
}
