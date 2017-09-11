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

func TestGetFeatureOfInterest(t *testing.T) {
	getAndAssertFeatureOfInterest("/v1.0/featuresofinterest(1)", t)
}

func TestGetFeatureOfInterestByObservation(t *testing.T) {
	getAndAssertFeatureOfInterest("/v1.0/observations(1)/featureofinterest", t)
}

func TestGetFeaturesOfInterest(t *testing.T) {
	getAndAssertFeaturesOfInterest("/v1.0/featuresofinterest", t)
}

func TestPostFeatureOfInterest(t *testing.T) {
	// arrange
	mockFeatureOfInterest := newMockFeatureOfInterest(1)

	// act
	r := request("POST", "/v1.0/featuresofinterest", mockFeatureOfInterest)

	// arrange
	parseAndAssertFeatureOfInterest(*mockFeatureOfInterest, r, http.StatusCreated, t)
}

func TestPutFeatureOfInterest(t *testing.T) {
	// arrange
	mockFeatureOfInterest := newMockFeatureOfInterest(1)
	mockFeatureOfInterest.Name = "patched"

	// act
	r := request("PUT", "/v1.0/featuresofinterest(1)", mockFeatureOfInterest)

	parseAndAssertFeatureOfInterest(*mockFeatureOfInterest, r, http.StatusOK, t)
}

func TestPatchFeatureOfInterest(t *testing.T) {
	// arrange
	mockFeatureOfInterest := newMockFeatureOfInterest(1)
	mockFeatureOfInterest.Name = "patched"

	// act
	r := request("PATCH", "/v1.0/featuresofinterest(1)", mockFeatureOfInterest)

	// assert
	parseAndAssertFeatureOfInterest(*mockFeatureOfInterest, r, http.StatusOK, t)
}

func TestDeleteFeatureOfInterest(t *testing.T) {
	r := request("DELETE", "/v1.0/featuresofinterest(1)", nil)

	// assert
	assertStatusCode(http.StatusOK, r, t)
}

func getAndAssertFeatureOfInterest(url string, t *testing.T) {
	r := request("GET", url, nil)
	parseAndAssertFeatureOfInterest(*newMockFeatureOfInterest(1), r, http.StatusOK, t)
}

func parseAndAssertFeatureOfInterest(created entities.FeatureOfInterest, r *http.Response, expectedStatusCode int, t *testing.T) {
	foi := entities.FeatureOfInterest{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &foi)

	assert.Nil(t, err)
	assertStatusCode(expectedStatusCode, r, t)
	assertFeatureOfInterest(created, foi, t)
}

func assertFeatureOfInterest(created, returned entities.FeatureOfInterest, t *testing.T) {
	assert.Equal(t, fmt.Sprintf("%v", created.ID), fmt.Sprintf("%v", returned.ID))
	assert.Equal(t, created.Name, returned.Name)
	assert.Equal(t, created.Description, returned.Description)
	assert.Equal(t, created.EncodingType, returned.EncodingType)
}

func getAndAssertFeaturesOfInterest(url string, t *testing.T) {
	// act
	r, _ := http.Get(getServer().URL + url)
	ar := entities.ArrayResponseFeaturesOfInterest{}
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &ar)

	// assert
	assert.Nil(t, err)
	assertStatusCode(http.StatusOK, r, t)
	assert.Equal(t, 2, ar.Count)

	for _, entity := range ar.Data {
		expected := newMockFeatureOfInterest(int(entity.ID.(float64)))
		assertFeatureOfInterest(*expected, *entity, t)
	}
}
