package config

import (
	"github.com/geodan/gost/sensorthings/models"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCreateEndPoints(t *testing.T) {
	//arrange
	endpoints := CreateEndPoints("http://test.com")

	//assert
	assert.Equal(t, 10, len(endpoints))
}

func TestCreateEndPointVersion(t *testing.T) {
	//arrange
	ve := CreateVersionEndpoint("http://test.com")

	//assert
	containsVersionPath := containsEndpoint("version", ve.Operations)
	assert.Equal(t, true, containsVersionPath, "Version endpoint needs to contain an endpoint containing the path Version")
}

func containsEndpoint(epName string, eps []models.EndpointOperation) bool {
	for _, o := range eps {
		if strings.Contains(o.Path, epName) {
			return true
		}
	}

	return false
}
