package entities

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var jsonThing = `{
		"description": "camping lantern",
		"properties": {
		"property1": "it’s waterproof"
		}
	}`

var jsonThingError = `{
		"desc": "camping lantern",
	}`

func TestMissingMandatoryParametersThing(t *testing.T) {
	//arrange
	thing := &Thing{}

	//act
	_, err := thing.ContainsMandatoryParams()

	//assert
	assert.NotNil(t, err, "Thing mandatory param description not filled in should have returned error")
	if len(err) > 0 {
		assert.Contains(t, fmt.Sprintf("%v", err[0]), "description")
	}
}

func TestMandatoryParametersExistThing(t *testing.T) {
	//arrange
	thing := &Thing{Description: "test"}

	//act
	_, err := thing.ContainsMandatoryParams()

	//assert
	assert.Nil(t, err, "All mandatory params are filled in shoud not have return an error")
}

func TestParseEntityResultOkThing(t *testing.T) {
	//arrange
	thing := &Thing{}

	//act
	err := thing.ParseEntity([]byte(jsonThing))

	//assert
	assert.Equal(t, err, nil, "Unable to parse json into thing")
}

func TestParseEntityResultNotOkThing(t *testing.T) {
	//arrange
	thing := &Thing{}

	//act
	err := thing.ParseEntity([]byte(jsonThingError))

	//assert
	assert.NotEqual(t, err, nil, "Thing parse from json should have failed")
}

func TestSetLinksThing(t *testing.T) {
	//arrange
	thing := &Thing{ID: id}

	//act
	thing.SetLinks(externalURL)

	//assert
	assert.Equal(t, thing.NavSelf, fmt.Sprintf("%s/v1.0/%s(%s)", externalURL, EntityLinkThings.ToString(), id), "Thing navself incorrect")
	assert.Equal(t, thing.NavDatastreams, fmt.Sprintf("../%s(%s)/%s", EntityLinkThings.ToString(), id, EntityLinkDatastreams.ToString()), "Thing NavDatastreams incorrect")
	assert.Equal(t, thing.NavLocations, fmt.Sprintf("../%s(%s)/%s", EntityLinkThings.ToString(), id, EntityLinkLocations.ToString()), "Thing NavLocations incorrect")
	assert.Equal(t, thing.NavHistoricalLocations, fmt.Sprintf("../%s(%s)/%s", EntityLinkThings.ToString(), id, EntityLinkHistoricalLocations.ToString()), "Thing NavHistoricalLocations incorrect")
}
