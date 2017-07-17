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

func TestGetThing(t *testing.T) {
	getAndAssertThing("/v1.0/things(1)", t)
}

func TestGetThingByDatastream(t *testing.T) {
	getAndAssertThing("/v1.0/datastreams(1)/thing", t)
}

func TestGetThingByHistoricalLocation(t *testing.T) {
	getAndAssertThing("/v1.0/historicallocations(1)/thing", t)
}

func TestGetThings(t *testing.T) {
	getAndAssertThings("/v1.0/things", t)
}

func TestGetThingsByLocation(t *testing.T) {
	getAndAssertThings("/v1.0/locations(1)/things", t)
}

func TestPostThing(t *testing.T) {
	// arrange
	mockThing := newMockThing(1)

	// act
	r := request("POST", "/v1.0/things", mockThing)

	// arrange
	parseAndAssertThing(*mockThing, r, http.StatusCreated, t)
}

func TestPutThing(t *testing.T) {
	// arrange
	mockThing := newMockThing(1)
	mockThing.Name = "patched"

	// act
	r := request("PUT", "/v1.0/things(1)", mockThing)

	parseAndAssertThing(*mockThing, r, http.StatusOK, t)
}

func TestPatchThing(t *testing.T) {
	// arrange
	mockThing := newMockThing(1)
	mockThing.Name = "patched"

	// act
	r := request("PATCH", "/v1.0/things(1)", mockThing)

	// assert
	parseAndAssertThing(*mockThing, r, http.StatusOK, t)
}

func TestDeleteThing(t *testing.T) {
	r := request("DELETE", "/v1.0/things(1)", nil)

	// assert
	assertStatusCode(http.StatusOK, r, t)
}

func getAndAssertThing(url string, t *testing.T) {
	r := request("GET", url, nil)
	parseAndAssertThing(*newMockThing(1), r, http.StatusOK, t)
}

func parseAndAssertThing(created entities.Thing, r *http.Response, expectedStatusCode int, t *testing.T) {
	thing := entities.Thing{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &thing)

	assert.Nil(t, err)
	assertStatusCode(expectedStatusCode, r, t)
	assertThing(created, thing, t)
}

func assertThing(created, returned entities.Thing, t *testing.T) {
	assert.Equal(t, fmt.Sprintf("%v", created.ID), fmt.Sprintf("%v", returned.ID))
	assert.Equal(t, created.Name, returned.Name)
	assert.Equal(t, created.Description, returned.Description)
	assert.Equal(t, created.Properties, returned.Properties)
}

func getAndAssertThings(url string, t *testing.T) {
	// act
	r, _ := http.Get(getServer().URL + url)
	ar := models.ArrayResponseThings{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &ar)

	// assert
	assert.Nil(t, err)
	assertStatusCode(http.StatusOK, r, t)
	assert.Equal(t, 2, ar.Count)

	for _, thing := range ar.Data {
		expected := newMockThing(int(thing.ID.(float64)))
		assertThing(*expected, *thing, t)
	}
}
