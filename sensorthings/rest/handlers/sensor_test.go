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

func TestGetSensor(t *testing.T) {
	getAndAssertSensor("/v1.0/sensors(1)", t)
}

func TestGetSensorByDatastream(t *testing.T) {
	getAndAssertSensor("/v1.0/datastreams(1)/sensor", t)
}

func TestGetSensors(t *testing.T) {
	getAndAssertSensors("/v1.0/sensors", t)
}

func TestPostSensor(t *testing.T) {
	// arrange
	mockSensor := newMockSensor(1)

	// act
	r := request("POST", "/v1.0/sensors", mockSensor)

	// arrange
	parseAndAssertSensor(*mockSensor, r, http.StatusCreated, t)
}

func TestPutSensor(t *testing.T) {
	// arrange
	mockSensor := newMockSensor(1)
	mockSensor.Name = "patched"

	// act
	r := request("PUT", "/v1.0/sensors(1)", mockSensor)

	parseAndAssertSensor(*mockSensor, r, http.StatusOK, t)
}

func TestPatchSensor(t *testing.T) {
	// arrange
	mockSensor := newMockSensor(1)
	mockSensor.Name = "patched"

	// act
	r := request("PATCH", "/v1.0/sensors(1)", mockSensor)

	// assert
	parseAndAssertSensor(*mockSensor, r, http.StatusOK, t)
}

func TestDeleteSensor(t *testing.T) {
	r := request("DELETE", "/v1.0/sensors(1)", nil)

	// assert
	assertStatusCode(http.StatusOK, r, t)
}

func getAndAssertSensor(url string, t *testing.T) {
	r := request("GET", url, nil)
	parseAndAssertSensor(*newMockSensor(1), r, http.StatusOK, t)
}

func parseAndAssertSensor(created entities.Sensor, r *http.Response, expectedStatusCode int, t *testing.T) {
	sensor := entities.Sensor{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &sensor)

	assert.Nil(t, err)
	assertStatusCode(expectedStatusCode, r, t)
	assertSensor(created, sensor, t)
}

func assertSensor(created, returned entities.Sensor, t *testing.T) {
	assert.Equal(t, fmt.Sprintf("%v", created.ID), fmt.Sprintf("%v", returned.ID))
	assert.Equal(t, created.Name, returned.Name)
	assert.Equal(t, created.Description, returned.Description)
	assert.Equal(t, created.EncodingType, returned.EncodingType)
	assert.Equal(t, created.Metadata, returned.Metadata)
}

func getAndAssertSensors(url string, t *testing.T) {
	// act
	r, _ := http.Get(getServer().URL + url)
	ar := models.ArrayResponseSensors{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &ar)

	// assert
	assert.Nil(t, err)
	assertStatusCode(http.StatusOK, r, t)
	assert.Equal(t, 2, ar.Count)

	for _, entity := range ar.Data {
		expected := newMockSensor(int(entity.ID.(float64)))
		assertSensor(*expected, *entity, t)
	}
}
