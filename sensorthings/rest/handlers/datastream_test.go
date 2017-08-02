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

func TestGetDatastream(t *testing.T) {
	getAndAssertDatastream("/v1.0/datastreams(1)", t)
}

func TestGetDatastreams(t *testing.T) {
	getAndAssertDatastreams("/v1.0/datastreams", t)
}

func TestGetDatastreamByObservation(t *testing.T) {
	getAndAssertDatastream("/v1.0/observations(1)/datastream", t)
}

func TestGetDatastreamsByThing(t *testing.T) {
	getAndAssertDatastreams("/v1.0/things(1)/datastreams", t)
}

func TestGetDatastreamsBySensor(t *testing.T) {
	getAndAssertDatastreams("/v1.0/sensors(1)/datastreams", t)
}

func TestGetDatastreamsByObservedProperty(t *testing.T) {
	getAndAssertDatastreams("/v1.0/observedproperties(1)/datastreams", t)
}

func TestPostDatastream(t *testing.T) {
	// arrange
	mockDatastream:= newMockDatastream(1)

	// act
	r := request("POST", "/v1.0/datastreams", mockDatastream)

	// arrange
	parseAndAssertDatastream(*mockDatastream, r, http.StatusCreated, t)
}

func TestPostDatastreamByThing(t *testing.T) {
	// arrange
	mockDatastream:= newMockDatastream(1)

	// act
	r := request("POST", "/v1.0/things(1)/datastreams", mockDatastream)

	// arrange
	parseAndAssertDatastream(*mockDatastream, r, http.StatusCreated, t)
}

func TestPutDatastream(t *testing.T) {
	// arrange
	mockDatastream := newMockDatastream(1)
	mockDatastream.Name = "patched"

	// act
	r := request("PUT", "/v1.0/datastreams(1)", mockDatastream)

	parseAndAssertDatastream(*mockDatastream, r, http.StatusOK, t)
}

func TestPatchDatastream(t *testing.T) {
	// arrange
	mockDatastream := newMockDatastream(1)
	mockDatastream.Name = "patched"

	// act
	r := request("PATCH", "/v1.0/datastreams(1)", mockDatastream)

	// assert
	parseAndAssertDatastream(*mockDatastream, r, http.StatusOK, t)
}

func TestDeleteDatastream(t *testing.T) {
	r := request("DELETE", "/v1.0/datastreams(1)", nil)

	// assert
	assertStatusCode(http.StatusOK, r, t)
}

func getAndAssertDatastream(url string, t *testing.T) {
	r := request("GET", url, nil)
	parseAndAssertDatastream(*newMockDatastream(1), r, http.StatusOK, t)
}

func parseAndAssertDatastream(created entities.Datastream, r *http.Response, expectedStatusCode int, t *testing.T) {
	d := entities.Datastream{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &d)

	assert.Nil(t, err)
	assertStatusCode(expectedStatusCode, r, t)
	assertDatastream(created, d, t)
}

func assertDatastream(created, returned entities.Datastream, t *testing.T) {
	assert.Equal(t, fmt.Sprintf("%v", created.ID), fmt.Sprintf("%v", returned.ID))
	assert.Equal(t, created.Name, returned.Name)
	assert.Equal(t, created.Description, returned.Description)
}

func getAndAssertDatastreams(url string, t *testing.T) {
	// act
	r, _ := http.Get(getServer().URL + url)
	ar := models.ArrayResponseDatastreams{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &ar)

	// assert
	assert.Nil(t, err)
	assertStatusCode(http.StatusOK, r, t)
	assert.Equal(t, 2, ar.Count)

	for _, entity := range ar.Data {
		expected := newMockDatastream(int(entity.ID.(float64)))
		assertDatastream(*expected, *entity, t)
	}
}
