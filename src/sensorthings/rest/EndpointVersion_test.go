package rest

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetEndPointVersion(t *testing.T) {
	// arrange
	ep := createVersionEndpoint("http://www.nu.nl")

	// assert
	assert.True(t, ep.URL == "http://www.nu.nl")
}
