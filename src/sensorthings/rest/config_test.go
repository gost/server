package rest

import (
	"github.com/geodan/gost/src/sensorthings/models"
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
	ve := createVersion("http://test.com")

	//assert
	containsVersionPath := containsEndpoint("Version", ve.Operations)
	assert.Equal(t, true, containsVersionPath, "Version endpoint needs to contain an endpoint containing the path Version")
}

func containsEndpoint(epName string, eps []models.EndpointOperation) bool {
	for _, o := range eps {
		if strings.Contains(o.Path, "Version") {
			return true
		}
	}

	return false
}
