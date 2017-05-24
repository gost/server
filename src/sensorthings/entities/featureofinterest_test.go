package entities

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetLinks(t *testing.T) {
	// arrange
	foi := &FeatureOfInterest{}
	foi.ID = "0"
	obs1 := &Observation{}
	obs2 := &Observation{}
	obs := []*Observation{obs1, obs2}
	foi.Observations = obs

	// act
	foi.SetAllLinks("http://www.nu.nl")
	propertynames := foi.GetPropertyNames()

	// assert
	assert.NotNil(t, foi.NavSelf, "NavSelf should be filled in")
	assert.True(t, len(propertynames) > 0)
}

func TestParseEntityFoi(t *testing.T) {
	// arrange
	foi := &FeatureOfInterest{}
	foi.ID = "0"
	dsjson, _ := json.Marshal(foi)

	// act
	foi.ParseEntity(dsjson)

	// assert
	assert.True(t, foi.ID == "0")
}

func TestParseEntityFoiFails(t *testing.T) {
	// arrange
	foi := &FeatureOfInterest{}
	dsjson, _ := json.Marshal("hoho")

	// act
	err := foi.ParseEntity(dsjson)

	// assert
	assert.True(t, err != nil)
}

func TestMissingMandatoryParametersFeatureOfInterest(t *testing.T) {
	//arrange
	featureofinterest := &FeatureOfInterest{}

	//act
	_, err := featureofinterest.ContainsMandatoryParams()

	//assert
	assert.NotNil(t, err, "FeatureOfInterest mandatory param description not filled in should have returned error")
	if len(err) > 0 {
		assert.Contains(t, fmt.Sprintf("%v", err[0]), "name")
	}
}

func TestMissingMandatoryParametersFeatureOfInterestSucceeds(t *testing.T) {
	//arrange
	featureofinterest := &FeatureOfInterest{}
	featureofinterest.Name = "name"
	featureofinterest.Description = "desc"
	featureofinterest.EncodingType = "type"
	featureofinterest.Feature = map[string]interface{}{"name": "location name 1"}

	//act
	_, err := featureofinterest.ContainsMandatoryParams()

	//assert
	assert.Nil(t, err, "FeatureOfInterest mandatory param description filled in should not have returned error")
}
