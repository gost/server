package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetLinksReturnsVariousLinks(t *testing.T) {
	// arrange
	datastream := &Datastream{}
	datastream.ID="0"
	// todo: what happens if there is no ID?

	// act
	datastream.SetLinks("http://www.nu.nl")

	// assert
	assert.NotNil(t, datastream.NavSelf, "NavSelf should be filled in")
	assert.NotNil(t, datastream.NavThing, "NavThing should be filled in")
	assert.NotNil(t, datastream.NavSensor, "NavSensor should be filled in")
	assert.NotNil(t, datastream.NavObservations, "NavObservations should be filled in")
	assert.NotNil(t, datastream.NavObservedProperty, "NavObservedProperty should be filled in")
}

func TestGetEntityType(t *testing.T){
	// arrange
	datastream := &Datastream{}

	// act
	entityType := datastream.GetEntityType()

	// assert
	assert.Equal(t, EntityTypeDatastream,entityType,"GetEntityType should be Datastream")
}

func TestContainsMandatoryParameters(t *testing.T){
	// arrange
	datastream := &Datastream{}

	// act
	contains, _ :=datastream.ContainsMandatoryParams()

	// assert
	assert.False(t,contains, "Datastream is expected not to have mandatory paramaters")
}

