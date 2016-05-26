package entities

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var jsonLocation = `{
    "description": "my backyard",
    "encodingType": "application/vnd.geo+json",
    "location": {
        "type": "Point",
        "coordinates": [-117.123,
        54.123]
    }
}`

var jsonLocationError = `{
    "description": "my backyard",
    "encodingType": "application/vnd.geo+json",
    "location": 10
}`

func TestMissingMandatoryParametersLocation(t *testing.T) {
	//arrange
	location := &Location{}

	//act
	_, err := location.ContainsMandatoryParams()

	//assert
	assert.NotNil(t, err, "Location mandatory param description not filled in should have returned error")
	if len(err) > 0 {
		assert.Contains(t, fmt.Sprintf("%v", err[0]), "description")
	}
}

func TestMandatoryParametersExistLocation(t *testing.T) {
	//arrange
	location := &Location{
		Description:  "test",
		EncodingType: "test",
		Location:     map[string]interface{}{"test": "test"},
	}

	//act
	_, err := location.ContainsMandatoryParams()

	//assert
	assert.Nil(t, err, "All mandatory params are filled in shoud not have returned an error")
}

func TestParseEntityResultOkLocation(t *testing.T) {
	//arrange
	location := &Location{}

	//act
	err := location.ParseEntity([]byte(jsonLocation))

	//assert
	assert.Equal(t, err, nil, "Unable to parse json into Location")
}

func TestParseEntityResultNotOkLocation(t *testing.T) {
	//arrange
	location := &Location{}

	//act
	err := location.ParseEntity([]byte(jsonLocationError))

	//assert
	assert.NotEqual(t, err, nil, "Location parse from json should have failed")
}

func TestSetLinksLocation(t *testing.T) {
	//arrange
	location := &Location{}
	location.ID = id

	//act
	location.SetLinks(externalURL)

	//assert
	assert.Equal(t, location.NavSelf, fmt.Sprintf("%s/v1.0/%s(%s)", externalURL, EntityLinkLocations.ToString(), id), "Location navself incorrect")
	assert.Equal(t, location.NavThings, fmt.Sprintf("%s/v1.0/%s(%s)/%s", externalURL, EntityLinkLocations.ToString(), id, EntityLinkThings.ToString()), "Location NavThings incorrect")
	assert.Equal(t, location.NavHistoricalLocations, fmt.Sprintf("%s/v1.0/%s(%s)/%s", externalURL, EntityLinkLocations.ToString(), id, EntityLinkHistoricalLocations.ToString()), "Location NavHistoricalLocations incorrect")
}

func TestGetSupportedEncodingLocation(t *testing.T) {
	//arrange
	location := &Location{}

	//act
	encodings := location.GetSupportedEncoding()
	_, ok := encodings[EncodingGeoJSON.Code]

	//assert
	assert.Equal(t, 1, len(encodings), "Location should support 1 encoding")
	assert.Equal(t, true, ok, "Location should support EncodingGeoJSON")
}
