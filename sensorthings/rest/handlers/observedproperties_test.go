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

func TestGetObservedProperty(t *testing.T) {
	getAndAssertObservedProperty("/v1.0/observedproperties(1)", t)
}

func TestGetObservedPropertyByDatastream(t *testing.T) {
	getAndAssertObservedProperty("/v1.0/datastreams(1)/observedproperty", t)
}

func TestGetObservedProperties(t *testing.T) {
	getAndAssertObservedProperties("/v1.0/observedproperties", t)
}

func TestPostObservedProperty(t *testing.T) {
	// arrange
	mockObservedProperty := newMockObservedProperty(1)

	// act
	r := request("POST", "/v1.0/observedproperties", mockObservedProperty)

	// arrange
	parseAndAssertObservedProperty(*mockObservedProperty, r, http.StatusCreated, t)
}

func TestPutObservedProperty(t *testing.T) {
	// arrange
	mockObservedProperty := newMockObservedProperty(1)
	mockObservedProperty.Name = "patched"

	// act
	r := request("PUT", "/v1.0/observedproperties(1)", mockObservedProperty)

	parseAndAssertObservedProperty(*mockObservedProperty, r, http.StatusOK, t)
}

func TestPatchObservedProperty(t *testing.T) {
	// arrange
	mockObservedProperty := newMockObservedProperty(1)
	mockObservedProperty.Name = "patched"

	// act
	r := request("PATCH", "/v1.0/observedproperties(1)", mockObservedProperty)

	// assert
	parseAndAssertObservedProperty(*mockObservedProperty, r, http.StatusOK, t)
}

func TestDeleteObservedProperty(t *testing.T) {
	r := request("DELETE", "/v1.0/observedproperties(1)", nil)

	// assert
	assertStatusCode(http.StatusOK, r, t)
}

func getAndAssertObservedProperty(url string, t *testing.T) {
	r := request("GET", url, nil)
	parseAndAssertObservedProperty(*newMockObservedProperty(1), r, http.StatusOK, t)
}

func parseAndAssertObservedProperty(created entities.ObservedProperty, r *http.Response, expectedStatusCode int, t *testing.T) {
	observedProperty := entities.ObservedProperty{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &observedProperty)

	assert.Nil(t, err)
	assertStatusCode(expectedStatusCode, r, t)
	assertObservedProperty(created, observedProperty, t)
}

func assertObservedProperty(created, returned entities.ObservedProperty, t *testing.T) {
	assert.Equal(t, fmt.Sprintf("%v", created.ID), fmt.Sprintf("%v", returned.ID))
	assert.Equal(t, created.Name, returned.Name)
	assert.Equal(t, created.Description, returned.Description)
}

func getAndAssertObservedProperties(url string, t *testing.T) {
	// act
	r, _ := http.Get(getServer().URL + url)
	ar := entities.ArrayResponseObservedProperty{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &ar)

	// assert
	assert.Nil(t, err)
	assertStatusCode(http.StatusOK, r, t)
	assert.Equal(t, 2, ar.Count)

	for _, entity := range ar.Data {
		expected := newMockObservedProperty(int(entity.ID.(float64)))
		assertObservedProperty(*expected, *entity, t)
	}
}
