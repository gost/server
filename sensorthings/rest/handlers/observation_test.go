package handlers

import (
	"encoding/json"
	"fmt"
	entities "github.com/gost/core"
	"github.com/gost/server/sensorthings/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetObservation(t *testing.T) {
	getAndAssertObservation("/v1.0/observations(1)", t)
}

func TestGetObservations(t *testing.T) {
	getAndAssertObservations("/v1.0/observations", t)
}

func TestGetObservationByDatastream(t *testing.T) {
	getAndAssertObservations("/v1.0/datastreams(1)/observations", t)
}

func TestGetObservationByFOI(t *testing.T) {
	getAndAssertObservations("/v1.0/featuresofinterest(1)/observations", t)
}

func TestPostObservation(t *testing.T) {
	// arrange
	mockObs := newMockObservation(1)

	// act
	r := request("POST", "/v1.0/observations", mockObs)

	// arrange
	parseAndAssertObservation(*mockObs, r, http.StatusCreated, t)
}

func TestPostObservationByDatastream(t *testing.T) {
	// arrange
	mockObs := newMockObservation(1)

	// act
	r := request("POST", "/v1.0/datastreams(1)/observations", mockObs)

	// arrange
	parseAndAssertObservation(*mockObs, r, http.StatusCreated, t)
}

func TestPutObservation(t *testing.T) {
	// arrange
	mockObs := newMockObservation(1)
	mockObs.PhenomenonTime = "2017-07-17T05:13:09.161Z"

	// act
	r := request("PUT", "/v1.0/observations(1)", mockObs)

	parseAndAssertObservation(*mockObs, r, http.StatusOK, t)
}

func TestPatchObservation(t *testing.T) {
	// arrange
	mockObs := newMockObservation(1)
	mockObs.PhenomenonTime = "2017-07-17T05:13:09.161Z"

	// act
	r := request("PATCH", "/v1.0/observations(1)", mockObs)

	// assert
	parseAndAssertObservation(*mockObs, r, http.StatusOK, t)
}

func TestDeleteObservations(t *testing.T) {
	r := request("DELETE", "/v1.0/observations(1)", nil)

	// assert
	assertStatusCode(http.StatusOK, r, t)
}

func getAndAssertObservation(url string, t *testing.T) {
	r := request("GET", url, nil)
	parseAndAssertObservation(*newMockObservation(1), r, http.StatusOK, t)
}

func parseAndAssertObservation(created entities.Observation, r *http.Response, expectedStatusCode int, t *testing.T) {
	obs := entities.Observation{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &obs)

	assert.Nil(t, err)
	assertStatusCode(expectedStatusCode, r, t)
	assertObservation(created, obs, t)
}

func assertObservation(created, returned entities.Observation, t *testing.T) {
	assert.Equal(t, fmt.Sprintf("%v", created.ID), fmt.Sprintf("%v", returned.ID))
	assert.Equal(t, fmt.Sprintf("%v", created.Result), fmt.Sprintf("%v", returned.Result))
	assert.Equal(t, created.PhenomenonTime, returned.PhenomenonTime)
}

func getAndAssertObservations(url string, t *testing.T) {
	// act
	r, _ := http.Get(getServer().URL + url)
	ar := models.ArrayResponseObservations{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &ar)

	// assert
	assert.Nil(t, err)
	assertStatusCode(http.StatusOK, r, t)
	assert.Equal(t, 2, ar.Count)

	for _, obs := range ar.Data {
		expected := newMockObservation(int(obs.ID.(float64)))
		assertObservation(*expected, *obs, t)
	}
}
