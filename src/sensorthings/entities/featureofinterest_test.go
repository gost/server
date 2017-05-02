package entities

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"encoding/json"
)

func TestSetLinks(t *testing.T) {
	// arrange
	foi := &FeatureOfInterest{}
	foi.ID = "0"

	// act
	foi.SetAllLinks("http://www.nu.nl")
	propertynames := foi.GetPropertyNames()
	supportedencoding := foi.GetSupportedEncoding()

	// assert
	assert.NotNil(t, foi.NavSelf, "NavSelf should be filled in")
	assert.True(t,len(propertynames)>0)
	assert.True(t,len(supportedencoding)>0)
}

func TestParseEntityFoi(t *testing.T){
	// arrange
	foi := &FeatureOfInterest{}
	foi.ID = "0"
	dsjson,_:=json.Marshal(foi)

	// act
	foi.ParseEntity(dsjson)

	// assert
	assert.True(t,foi.ID == "0")
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
