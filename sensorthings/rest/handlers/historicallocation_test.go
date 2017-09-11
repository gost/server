package handlers

import (
	"encoding/json"
	"fmt"
	entities "github.com/gost/core"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetHistoricalLocation(t *testing.T) {
	getAndAssertHistoricalLocation("/v1.0/historicallocations(1)", t)
}

func TestGetHistoricalLocationsByThing(t *testing.T) {
	getAndAssertHistoricalLocations("/v1.0/things(1)/historicallocations", t)
}

func TestGetHistoricalLocationsByLocation(t *testing.T) {
	getAndAssertHistoricalLocations("/v1.0/locations(1)/historicallocations", t)
}

func TestGetHistoricalLocations(t *testing.T) {
	getAndAssertHistoricalLocations("/v1.0/historicallocations", t)
}

func TestPostHistoricalLocation(t *testing.T) {
	// arrange
	mockHistoricalLocation := newMockHistoricalLocation(1)

	// act
	r := request("POST", "/v1.0/historicallocations", mockHistoricalLocation)

	// arrange
	parseAndAssertHistoricalLocation(*mockHistoricalLocation, r, http.StatusCreated, t)
}

func TestPutHistoricalLocation(t *testing.T) {
	// arrange
	mockHistoricalLocation := newMockHistoricalLocation(1)
	mockHistoricalLocation.Time = "2017-07-17T07:04:09.222Z"

	// act
	r := request("PUT", "/v1.0/historicallocations(1)", mockHistoricalLocation)

	parseAndAssertHistoricalLocation(*mockHistoricalLocation, r, http.StatusOK, t)
}

func TestPatchHistoricalLocation(t *testing.T) {
	// arrange
	mockHistoricalLocation := newMockHistoricalLocation(1)
	mockHistoricalLocation.Time = "2017-07-17T07:04:09.222Z"

	// act
	r := request("PATCH", "/v1.0/historicallocations(1)", mockHistoricalLocation)

	// assert
	parseAndAssertHistoricalLocation(*mockHistoricalLocation, r, http.StatusOK, t)
}

func TestDeleteHistoricalLocation(t *testing.T) {
	r := request("DELETE", "/v1.0/historicallocations(1)", nil)

	// assert
	assertStatusCode(http.StatusOK, r, t)
}

func getAndAssertHistoricalLocation(url string, t *testing.T) {
	r := request("GET", url, nil)
	parseAndAssertHistoricalLocation(*newMockHistoricalLocation(1), r, http.StatusOK, t)
}

func parseAndAssertHistoricalLocation(created entities.HistoricalLocation, r *http.Response, expectedStatusCode int, t *testing.T) {
	historicalLocation := entities.HistoricalLocation{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &historicalLocation)

	assert.Nil(t, err)
	assertStatusCode(expectedStatusCode, r, t)
	assertHistoricalLocation(created, historicalLocation, t)
}

func assertHistoricalLocation(created, returned entities.HistoricalLocation, t *testing.T) {
	assert.Equal(t, fmt.Sprintf("%v", created.ID), fmt.Sprintf("%v", returned.ID))
	assert.Equal(t, created.Time, returned.Time)
}

func getAndAssertHistoricalLocations(url string, t *testing.T) {
	// act
	r, _ := http.Get(getServer().URL + url)
	ar := entities.ArrayResponseHistoricalLocations{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &ar)

	// assert
	assert.Nil(t, err)
	assertStatusCode(http.StatusOK, r, t)
	assert.Equal(t, 2, ar.Count)

	for _, entity := range ar.Data {
		expected := newMockHistoricalLocation(int(entity.ID.(float64)))
		assertHistoricalLocation(*expected, *entity, t)
	}
}
