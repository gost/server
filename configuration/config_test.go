package configuration

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetExternalUri(t *testing.T) {
    var testurl = "http://test.com/"
    // arrange
	cfg := Config{}
	cfg.Server.ExternalURI = testurl
    // act
    var uri = cfg.GetExternalServerURI()
    // assert
	assert.Equal(t, uri, "http://test.com", "Trailing slash not removed by GetExternalServerUri")
}

func TestGetInternalUri(t *testing.T) {
	// arrange
	cfg := Config{}
	cfg.Server.Host = "localhost"
	cfg.Server.Port = 8080
	// act
	var uri = cfg.GetInternalServerURI()
	// assert
	assert.Equal(t, "localhost:8080", uri, "Internal server uri not constructed correctly based on config server host and port")
}

