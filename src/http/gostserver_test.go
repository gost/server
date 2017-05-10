package http

import (
	"github.com/geodan/gost/src/configuration"
	"github.com/geodan/gost/src/database/postgis"
	"github.com/geodan/gost/src/mqtt"
	api "github.com/geodan/gost/src/sensorthings/api"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateServer(t *testing.T) {
	// act
	cfg := configuration.Config{}
	mqttServer := mqtt.CreateMQTTClient(configuration.MQTTConfig{})
	database := postgis.NewDatabase("", 123, "", "", "", "", false, 50, 100, 200)
	stAPI := api.NewAPI(database, cfg, mqttServer)
	server := CreateServer("localhost", 8000, &stAPI, false, "", "")
	server.Stop()

	// assert
	assert.NotNil(t, server)
}
