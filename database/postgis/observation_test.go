package postgis

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/geodan/gost/sensorthings/entities"
)

func TestObservationParamFactory(t *testing.T) {
	// arrange
	phenomenonTime := "2015-03-06T00:00:00.000Z"
	resultTime := "2015-03-06T00:00:00.000Z"
	validTime := "2015-03-06T00:00:00.000Z"

	values := map[string]interface{}{
		"observation_id":             4,
		"observation_phenomenontime": phenomenonTime,
		"observation_result": "!0.5",
		"observation_resulttime": resultTime,
		"observation_resultquality": "goed",
		"observation_validtime": validTime,
		//"observation_parameters": "test",
	}

	// act
	entity, err := observationParamFactory(values)
	entitytype := entity.GetEntityType()
	// todo: how to get the observation??

	// assert
	assert.True(t, entity != nil)
	// entities..
	assert.True(t, err == nil)
	assert.True(t,entity.GetID() == 4)
	assert.True(t,entitytype == entities.EntityTypeObservation)
	// assert.True(t,*observation.ResultTime == resultTime)
}
