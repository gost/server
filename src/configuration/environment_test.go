package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentVariavels(t *testing.T) {
	// arrange
	conf := Config{}

	// act
	SetEnvironmentVariables(&conf)

	// assert
	assert.NotNil(t, conf, "Configuration should not be nil")
}
