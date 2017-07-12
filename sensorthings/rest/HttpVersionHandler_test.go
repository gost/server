package rest

import (
	"net/http"

	"encoding/json"
	"github.com/geodan/gost/configuration"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestVersionResponse(t *testing.T) {
	// arrange
	a := NewMockAPI()
	rr := httptest.NewRecorder()

	// act
	HandleVersion(rr, nil, nil, &a)
	version := models.VersionInfo{}
	json.Unmarshal(rr.Body.Bytes(), &version)
	decoder := json.NewDecoder(rr.Body)
	err := decoder.Decode(&version)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, configuration.SensorThingsAPIVersion, version.APIVersion.Version)
	assert.Equal(t, configuration.ServerVersion, version.GostServerVersion.Version)
}
