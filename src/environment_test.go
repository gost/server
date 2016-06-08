package main

import (
	"testing"

	"github.com/geodan/gost/src/configuration"
	"github.com/stretchr/testify/assert"
)

func TestEnvironmentVariavels(t *testing.T) {
	// arrange
	conf := configuration.Config{}

	// act
	SetEnvironmentVariables(&conf)

	// assert
	assert.NotNil(t, conf, "Configuration should not be nil")
}
