package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetThing(t *testing.T) {
	testGetAndAssertThing("/v1.0/things(1)", t)
}

func TestGetThingByDatastream(t *testing.T) {
	testGetAndAssertThing("/v1.0/datastreams(1)/thing", t)
}

func TestGetThingByHistoricalLocation(t *testing.T) {
	testGetAndAssertThing("/v1.0/historicallocations(1)/thing", t)
}

func TestGetThings(t *testing.T) {
	testThings("/v1.0/things", t)
}

func TestGetThingsByLocation(t *testing.T) {
	testThings("/v1.0/locations(1)/things", t)
}

func TestPostThing(t *testing.T) {
	// arrange
	mockThing := NewMockThing(1)
	b, _ := json.Marshal(mockThing)
	client := &http.Client{}
	request, _ := http.NewRequest("POST", getServer().URL+"/v1.0/things", bytes.NewReader(b))

	// act
	r, _ := client.Do(request)

	// arrange
	assertThing(*mockThing, r, http.StatusCreated, t)
}

func TestPutThing(t *testing.T) {
	// arrange
	mockThing := NewMockThing(1)
	mockThing.Name = "patched"
	b, _ := json.Marshal(mockThing)
	client := &http.Client{}
	request, _ := http.NewRequest("PUT", getServer().URL+"/v1.0/things(1)", bytes.NewReader(b))

	// act
	r, _ := client.Do(request)

	assertThing(*mockThing, r, http.StatusOK, t)
}

func TestPatchThing(t *testing.T) {
	// arrange
	mockThing := NewMockThing(1)
	mockThing.Name = "patched"
	b, _ := json.Marshal(mockThing)
	client := &http.Client{}
	request, _ := http.NewRequest("PATCH", getServer().URL+"/v1.0/things(1)", bytes.NewReader(b))

	// act
	r, _ := client.Do(request)

	assertThing(*mockThing, r, http.StatusOK, t)
}

func TestDeleteThing(t *testing.T) {
	// arrange
	client := &http.Client{}
	request, _ := http.NewRequest("DELETE", getServer().URL+"/v1.0/things(1)", nil)

	// act
	r, _ := client.Do(request)

	// assert
	assertStatusCode(http.StatusOK, r, t)
}

func testGetAndAssertThing(url string, t *testing.T) {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", getServer().URL+url, nil)
	r, _ := client.Do(request)
	assertThing(*NewMockThing(1), r, http.StatusOK, t)
}

func assertThing(created entities.Thing, r *http.Response, expectedStatusCode int, t *testing.T) {
	thing := entities.Thing{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &thing)

	assert.Nil(t, err)
	assertStatusCode(expectedStatusCode, r, t)
	assert.Equal(t, fmt.Sprintf("%v", created.ID), fmt.Sprintf("%v", thing.ID))
	assert.Equal(t, created.Name, thing.Name)
	assert.Equal(t, created.Description, thing.Description)
	assert.Equal(t, created.Properties, thing.Properties)
}

func testThings(url string, t *testing.T) {
	// act
	r, _ := http.Get(getServer().URL + url)
	ar := models.ArrayResponseThings{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &ar)

	// assert
	assert.Nil(t, err)
	assertStatusCode(http.StatusOK, r, t)
	assert.Equal(t, 2, ar.Count)
	assert.Equal(t, fmt.Sprintf("%v", ar.Data[0].ID), fmt.Sprintf("%v", 1))
	assert.Equal(t, fmt.Sprintf("%v", ar.Data[1].ID), fmt.Sprintf("%v", 2))
}

func assertStatusCode(expectedStatusCode int, r *http.Response, t *testing.T) {
	assert.Equal(t, expectedStatusCode, r.StatusCode)
}
