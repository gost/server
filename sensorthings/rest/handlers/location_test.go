package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gost/server/sensorthings/entities"
	"github.com/gost/server/sensorthings/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetLocation(t *testing.T) {
	getAndAssertLocation("/v1.0/locations(1)", t)
}

func TestGetLocationsByThing(t *testing.T) {
	getAndAssertLocations("/v1.0/things(1)/locations", t)
}

func TestHandleGetLocationsByHistoricalLocation(t *testing.T) {
	getAndAssertLocations("/v1.0/historicallocations(1)/locations", t)
}

func TestGetLocations(t *testing.T) {
	getAndAssertLocations("/v1.0/locations", t)
}

func TestPostLocation(t *testing.T) {
	// arrange
	mockLocation := newMockLocation(1)

	// act
	r := request("POST", "/v1.0/locations", mockLocation)

	// arrange
	parseAndAssertLocation(*mockLocation, r, http.StatusCreated, t)
}

func TestPostLocationByThing(t *testing.T){
	// arrange
	mockLocation := newMockLocation(1)

	// act
	r := request("POST", "/v1.0/things(1)/locations", mockLocation)

	// arrange
	parseAndAssertLocation(*mockLocation, r, http.StatusCreated, t)
}

func TestPutLocation(t *testing.T) {
	// arrange
	mockLocation := newMockLocation(1)
	mockLocation.Name = "patched"

	// act
	r := request("PUT", "/v1.0/locations(1)", mockLocation)

	parseAndAssertLocation(*mockLocation, r, http.StatusOK, t)
}

func TestPatchLocation(t *testing.T) {
	// arrange
	mockLocation := newMockLocation(1)
	mockLocation.Name = "patched"

	// act
	r := request("PATCH", "/v1.0/locations(1)", mockLocation)

	// assert
	parseAndAssertLocation(*mockLocation, r, http.StatusOK, t)
}

func TestDeleteLocation(t *testing.T) {
	r := request("DELETE", "/v1.0/locations(1)", nil)

	// assert
	assertStatusCode(http.StatusOK, r, t)
}

func getAndAssertLocation(url string, t *testing.T) {
	r := request("GET", url, nil)
	parseAndAssertLocation(*newMockLocation(1), r, http.StatusOK, t)
}

func parseAndAssertLocation(created entities.Location, r *http.Response, expectedStatusCode int, t *testing.T) {
	location := entities.Location{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &location)

	assert.Nil(t, err)
	assertStatusCode(expectedStatusCode, r, t)
	assertLocation(created, location, t)
}

func assertLocation(created, returned entities.Location, t *testing.T) {
	assert.Equal(t, fmt.Sprintf("%v", created.ID), fmt.Sprintf("%v", returned.ID))
	assert.Equal(t, created.Name, returned.Name)
	assert.Equal(t, created.Description, returned.Description)
	assert.Equal(t, created.EncodingType, returned.EncodingType)
	assert.Equal(t, created.Location, returned.Location)
}

func getAndAssertLocations(url string, t *testing.T) {
	// act
	r, _ := http.Get(getServer().URL + url)
	ar := models.ArrayResponseLocations{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &ar)

	// assert
	assert.Nil(t, err)
	assertStatusCode(http.StatusOK, r, t)
	assert.Equal(t, 2, ar.Count)

	for _, entity := range ar.Data {
		expected := newMockLocation(int(entity.ID.(float64)))
		assertLocation(*expected, *entity, t)
	}
}
