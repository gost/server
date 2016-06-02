package http

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/geodan/gost/src/sensorthings/api"
	"github.com/geodan/gost/src/database/postgis"
	"github.com/geodan/gost/src/mqtt"
	"github.com/geodan/gost/src/configuration"
)

// TestHttp starts
func TestCreateRouter(t *testing.T) {
	// arrange
	cfg := configuration.Config{}
	mqttServer := mqtt.CreateMQTTClient(configuration.MQTTConfig{})
	database := postgis.NewDatabase("", 123, "", "", "", "", false, 50, 100)
	a := api.NewAPI(database, cfg, mqttServer)
	
	// act
	router := CreateRouter(&a)

	// assert
	assert.NotNil(t, router, "Router should be created")
}
