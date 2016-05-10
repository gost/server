package api

import (
	"github.com/geodan/gost/configuration"
	"github.com/geodan/gost/database/postgis"
	"github.com/geodan/gost/mqtt"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateApi(t *testing.T) {
	// arrange
	cfg := configuration.Config{}
	mqttServer := mqtt.CreateMQTTClient(configuration.MQTTConfig{})
	database := postgis.NewDatabase("", 123, "", "", "", "", false)
	stAPI := NewAPI(database, cfg, mqttServer)

	// act
	ep := stAPI.GetEndpoints()

	// assert
	assert.NotNil(t, ep)
}
